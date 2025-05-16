export default [
    {
        name: "sflowlog_list",
        path: "/sflowlog/:id",
        component: () => import('@/pages/sflow/Log.vue'),
        meta: {
            auth: 'sflowlog.view',
        }
    },
]