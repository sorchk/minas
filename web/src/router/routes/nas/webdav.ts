export default [
    {
        name: "webdav_list",
        path: "/nas/webdav",
        component: () => import('@/pages/nas/webdav/List.vue'),
        meta: {
            auth: 'webdav.view',
        }
    },
    {
        name: "webdav_new",
        path: "/nas/webdav/new",
        component: () => import('@/pages/nas/webdav/Edit.vue'),
        meta: {
            auth: 'webdav.edit',
        }
    },
    {
        name: "webdav_detail",
        path: "/nas/webdav/:id",
        component: () => import('@/pages/nas/webdav/View.vue'),
        meta: {
            auth: 'webdav.view',
        }
    },
    {
        name: "webdav_edit",
        path: "/nas/webdav/:id/edit",
        component: () => import('@/pages/nas/webdav/Edit.vue'),
        meta: {
            auth: 'webdav.edit',
        }
    }
]