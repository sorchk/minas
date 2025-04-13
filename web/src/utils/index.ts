import { useBreakpoint, useMemo } from 'vooks'
import { isRef, Ref } from 'vue'

export function useIsMobile() {
  const breakpointRef = useBreakpoint()
  return useMemo(() => {
    return breakpointRef.value === 'xs'
  })
}

export function useIsTablet() {
  const breakpointRef = useBreakpoint()
  return useMemo(() => {
    return breakpointRef.value === 's'
  })
}

/**
 * Wrapper for setTimeout
 * 
 * @param callback 
 * @param interval delay for first call (unit: milliseconds) 
 * @param delay delay for first call (unit: milliseconds)
 * @returns 
 */
export function useTimer(callback: () => void, interval: number | Ref<number>, delay: number = 0): () => void {
  var id: NodeJS.Timeout
  const fn = () => {
    callback()
    id = setTimeout(fn, isRef(interval) ? interval.value : interval)
  }
  const stop = () => clearTimeout(id)

  if (delay === 0) {
    fn()
  } else {
    id = setTimeout(fn, delay)
  }

  return stop
}

export function isEmpty(...arrs: (any[] | undefined)[]): boolean {
  return arrs.every(arr => !arr || !arr.length)
}

export function toTitle(s: string): string {
  return s ? s[0].toUpperCase() + s.substring(1) : s
}

export function guid() {
  return s4() + s4() + s4() + s4() + s4() + s4() + s4() + s4()
}

function s4() {
  return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1);
}


/**
 * 对象深克隆
 * @param obj 源对象
 * @returns 克隆后的对象
 */
export function deepClone(obj: any) {
  let newObj: any;
  try {
    newObj = obj.push ? [] : {};
  } catch (error) {
    newObj = {};
  }
  for (let attr in obj) {
    if (obj[attr] && typeof obj[attr] === 'object') {
      newObj[attr] = deepClone(obj[attr]);
    } else {
      newObj[attr] = obj[attr];
    }
  }
  return newObj;
}
export const copyText = (text: string) => {
  try {
    navigator.clipboard.writeText(text).then(() => {
      window.message.success('复制成功');
    });
  } catch (err) {
  }
}