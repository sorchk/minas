export default [
    {
        name: "externalNas_list",
        path: "/nas/external",
        component: () => import('@/pages/nas/external/List.vue'),
        meta: {
            auth: 'externalNas.view',
        }
    },
    {
        name: "externalNas_new",
        path: "/nas/external/new",
        component: () => import('@/pages/nas/external/Edit.vue'),
        meta: {
            auth: 'externalNas.edit',
        }
    },
    {
        name: "externalNas_detail",
        path: "/nas/external/:id",
        component: () => import('@/pages/nas/external/View.vue'),
        meta: {
            auth: 'externalNas.view',
        }
    },
    {
        name: "externalNas_edit",
        path: "/nas/external/:id/edit",
        component: () => import('@/pages/nas/external/Edit.vue'),
        meta: {
            auth: 'externalNas.edit',
        }
    }
]