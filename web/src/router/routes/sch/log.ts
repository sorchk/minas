export default [
    {
        name: "schlog_list",
        path: "/schlog/:id",
        component: () => import('@/pages/schtask/Log.vue'),
        meta: {
            auth: 'schlog.view',
        }
    },
]