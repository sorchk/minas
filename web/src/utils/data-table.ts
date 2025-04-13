import { isRef, onMounted, reactive } from "vue"
import { t } from "@/locales";

export function useDataTable(loader: Function, filter: Object | Function, autoFetch: boolean = true) {
    const state = reactive({
        loading: false,
        data: [],
    })
    const pagination = reactive({
        page: 1,
        pageCount: 1,
        pageSize: 15,
        itemCount: 0,
        showSizePicker: true,
        pageSizes: [10, 15, 20, 50, 100, 200, 500],
        prefix({ itemCount }: any) {
            return t('texts.records', { total: itemCount } as any, itemCount)
        }
    })
    const fetchData = async function (page: number = 1) {
        state.data = [];
        state.loading = true;
        try {
            let args = typeof filter === 'function' ? filter() : filter
            args = isRef(args) ? args.value : args
            let r = await loader({
                ...args,
                page: page,
                size: pagination.pageSize,
            });
            state.data = r.data || [];
            pagination.itemCount = r.total || 0
            pagination.page = page
            pagination.pageCount = Math.ceil(pagination.itemCount / pagination.pageSize)
        } finally {
            state.loading = false;
        }
    }
    const changePageSize = function (size: number) {
        pagination.page = 1
        pagination.pageSize = size
        fetchData()
    }

    if (autoFetch) {
        onMounted(fetchData)
    }

    return { state, pagination, fetchData, changePageSize }
}
