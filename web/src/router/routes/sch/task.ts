export default [
    {
        name: "schtask_list",
        path: "/schtask",
        component: () => import('@/pages/schtask/List.vue'),
        meta: {
            auth: 'schtask.view',
        }
    },
    {
        name: "schtask_new",
        path: "/schtask/new",
        component: () => import('@/pages/schtask/Edit.vue'),
        meta: {
            auth: 'schtask.edit',
        }
    },
    {
        name: "schtask_detail",
        path: "/schtask/:id",
        component: () => import('@/pages/schtask/View.vue'),
        meta: {
            auth: 'schtask.view',
        }
    },
    {
        name: "schtask_edit",
        path: "/schtask/:id/edit",
        component: () => import('@/pages/schtask/Edit.vue'),
        meta: {
            auth: 'schtask.edit',
        }
    }
]