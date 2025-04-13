package logger

import (
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Writer 日志写入器基础结构体
// 负责将日志内容写入文件并管理日志滚动
type Writer struct {
	m             Manager       // 日志管理器，处理日志文件滚动和生命周期
	file          *os.File      // 当前打开的日志文件
	absPath       string        // 日志文件的绝对路径
	fire          chan string   // 接收文件滚动信号的通道
	cf            *LoggerConfig // 日志配置信息
	rollingfilech chan string   // 日志文件滚动队列，用于管理历史日志文件
}

// AsynchronousWriter 异步日志写入器
// 继承自Writer，提供异步写入日志的能力
type AsynchronousWriter struct {
	Writer                 // 嵌入基础写入器
	ctx     chan int       // 上下文控制通道，用于终止异步写入循环
	queue   chan []byte    // 日志数据缓冲队列
	errChan chan error     // 错误通道，用于传递写入过程中的错误
	closed  int32          // 关闭标志，使用原子操作保证线程安全
	wg      sync.WaitGroup // 等待组，确保所有goroutine完成
}

// Close 关闭异步写入器
// 会刷新剩余的日志数据，关闭文件并释放资源
// 返回值: 如果有错误发生返回error，否则返回nil
func (w *AsynchronousWriter) Close() error {
	// 使用原子操作确保Close方法只被执行一次
	if atomic.CompareAndSwapInt32(&w.closed, 0, 1) {
		close(w.ctx) // 关闭上下文通道，通知写入循环退出
		w.onClose()  // 处理关闭前的清理工作

		// 安全关闭日志管理器
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("%v", r)
				}
			}()
			w.m.Close()
		}()
		return w.file.Close() // 关闭日志文件
	}
	return ErrClosed // 已经关闭，返回相应错误
}

// onClose 处理关闭前的清理工作
// 将队列中剩余的日志数据写入文件
func (w *AsynchronousWriter) onClose() {
	var err error
	for {
		select {
		case b := <-w.queue: // 从队列中获取日志数据
			if _, err = w.file.Write(b); err != nil {
				select {
				case w.errChan <- err: // 尝试发送错误
				default:
					_asyncBufferPool.Put(&b) // 回收缓冲区
					return
				}
			}
			_asyncBufferPool.Put(&b) // 回收缓冲区
		default:
			return // 队列为空时退出
		}
	}
}

// _asyncBufferPool 异步写入器的缓冲区对象池
// 通过复用缓冲区减少内存分配和垃圾回收压力
var _asyncBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, BufferSize)
	},
}

// NewWriterFromConfig 根据配置创建新的滚动写入器
// 参数 c: 日志配置
// 返回值: 滚动日志写入器和可能的错误
func NewWriterFromConfig(c *LoggerConfig) (RollingWriter, error) {
	// 验证配置参数
	if c.LogPath == "" || c.FileName == "" {
		return nil, ErrInvalidArgument
	}

	// 创建日志目录
	if err := os.MkdirAll(c.LogPath, 0700); err != nil {
		return nil, err
	}

	// 打开或创建日志文件
	filepath := FilePath(c)
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// 重定向标准错误到日志文件
	if err := dupWrite(file); err != nil {
		return nil, err
	}

	// 创建日志管理器
	mng, err := NewManager(c)
	if err != nil {
		return nil, err
	}

	var rollingWriter RollingWriter
	// 初始化基础写入器
	writer := Writer{
		m:       mng,
		file:    file,
		absPath: filepath,
		fire:    mng.Fire(),
		cf:      c,
	}

	// 处理日志文件滚动和历史日志管理
	if c.MaxRemain > 0 {
		writer.rollingfilech = make(chan string, c.MaxRemain)
		// 读取当前日志目录
		dir, err := os.ReadDir(c.LogPath)
		if err != nil {
			mng.Close()
			return nil, err
		}

		// 查找符合命名规则的日志文件
		files := make([]string, 0, 10)
		for _, fi := range dir {
			if fi.IsDir() {
				continue
			}

			fileName := c.FileName
			if strings.Contains(fi.Name(), fileName) && strings.Contains(fi.Name(), c.LogSuffix) {
				start := strings.Index(fi.Name(), "-")
				end := strings.Index(fi.Name(), c.LogSuffix)
				name := fi.Name()
				if start > 0 && end > 0 {
					_, err := time.Parse(c.TimeTagFormat, name[start+1:end])
					if err == nil {
						files = append(files, fi.Name())
					}
				}
			}
		}

		// 按时间排序日志文件
		sort.Slice(files, func(i, j int) bool {
			t1Start := strings.Index(files[i], "-")
			t1End := strings.Index(files[i], c.LogSuffix)
			t2Start := strings.Index(files[i], "-")
			t2End := strings.Index(files[i], c.LogSuffix)
			t1, _ := time.Parse(c.TimeTagFormat, files[i][t1Start+1:t1End])
			t2, _ := time.Parse(c.TimeTagFormat, files[j][t2Start+1:t2End])
			return t1.Before(t2)
		})

		// 将日志文件放入滚动队列
		for _, file := range files {
		retry:
			select {
			case writer.rollingfilech <- path.Join(c.LogPath, file):
			default:
				writer.DoRemove() // 如果队列满了，移除最旧的日志文件
				goto retry
			}
		}
	}

	// 创建异步写入器
	wr := &AsynchronousWriter{
		ctx:     make(chan int),
		queue:   make(chan []byte, QueueSize),
		errChan: make(chan error, QueueSize),
		wg:      sync.WaitGroup{},
		closed:  0,
		Writer:  writer,
	}

	// 启动异步写入goroutine
	wr.wg.Add(1)
	go wr.writer()
	wr.wg.Wait()
	rollingWriter = wr

	return rollingWriter, nil
}

// writer 异步写入器的后台处理goroutine
// 负责处理日志文件滚动和实际写入操作
func (w *AsynchronousWriter) writer() {
	var err error
	w.wg.Done()
	for {
		select {
		case filename := <-w.fire: // 收到文件滚动信号
			if err = w.Reopen(filename); err != nil && len(w.errChan) < cap(w.errChan) {
				w.errChan <- err
			}
		case b := <-w.queue: // 收到日志数据
			if _, err = w.file.Write(b); err != nil && len(w.errChan) < cap(w.errChan) {
				w.errChan <- err
			}
			_asyncBufferPool.Put(&b) // 回收缓冲区
		case <-w.ctx: // 收到关闭信号
			return
		}
	}
}

// DoRemove 移除一个历史日志文件
// 从滚动队列中取出最旧的文件并删除
func (w *Writer) DoRemove() {
	file := <-w.rollingfilech
	if err := os.Remove(file); err != nil {
		log.Println("error in remove log file", file, err)
	}
}

// Write 实现io.Writer接口的写入方法
// 参数 b: 要写入的数据
// 返回值: 写入的字节数和可能的错误
func (w *Writer) Write(b []byte) (int, error) {
	// 处理可能存在的文件滚动信号
	var ok = false
	for !ok {
		select {
		case filename := <-w.fire:
			if err := w.Reopen(filename); err != nil {
				return 0, err
			}
		default:
			ok = true
		}
	}

	// 使用原子操作获取当前文件指针
	fp := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&w.file)))
	file := (*os.File)(fp)
	return file.Write(b)
}

// Reopen 重新打开日志文件
// 完成日志文件滚动操作
// 参数 file: 新日志文件名
// 返回值: 如果有错误发生返回error，否则返回nil
func (w *Writer) Reopen(file string) error {
	// 获取当前日志文件信息
	fileInfo, err := w.file.Stat()
	if err != nil {
		return err
	}

	// 如果文件为空，不需要滚动
	if fileInfo.Size() == 0 {
		return nil
	}

	// 关闭当前文件并重命名
	w.file.Close()
	if err := os.Rename(w.absPath, file); err != nil {
		return err
	}

	// 创建新的日志文件
	newFile, err := os.OpenFile(w.absPath, DefaultFileFlag, DefaultFileMode)
	if err != nil {
		return err
	}

	w.file = newFile

	// 处理历史日志文件
	go func() {
		if w.cf.MaxRemain > 0 {
		retry:
			select {
			case w.rollingfilech <- file: // 将滚动出的文件放入队列
			default:
				w.DoRemove() // 如果队列满了，移除最旧的日志文件
				goto retry
			}
		}
	}()
	return nil
}
