export default [
    {
        name: "projectdir_list",
        path: "/basic/projectdir",
        component: () => import('@/pages/basic/projectdir/List.vue'),
        meta: {
            auth: 'projectdir.view',
        }
    },
    {
        name: "projectdir_new",
        path: "/basic/projectdir/new",
        component: () => import('@/pages/basic/projectdir/Edit.vue'),
        meta: {
            auth: 'projectdir.edit',
        }
    },
    {
        name: "projectdir_detail",
        path: "/basic/projectdir/:id",
        component: () => import('@/pages/basic/projectdir/View.vue'),
        meta: {
            auth: 'projectdir.view',
        }
    },
    {
        name: "projectdir_edit",
        path: "/basic/projectdir/:id/edit",
        component: () => import('@/pages/basic/projectdir/Edit.vue'),
        meta: {
            auth: 'projectdir.edit',
        }
    }
]