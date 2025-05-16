export default [
    {
        name: "sflow_list",
        path: "/sflow",
        component: () => import('@/pages/sflow/List.vue'),
        meta: {
            auth: 'sflow.view',
        }
    },
    {
        name: "sflow_new",
        path: "/sflow/new",
        component: () => import('@/pages/sflow/Edit.vue'),
        meta: {
            auth: 'sflow.edit',
        }
    },
    {
        name: "sflow_detail",
        path: "/sflow/:id",
        component: () => import('@/pages/sflow/View.vue'),
        meta: {
            auth: 'sflow.view',
        }
    },
    {
        name: "sflow_edit",
        path: "/sflow/:id/edit",
        component: () => import('@/pages/sflow/Edit.vue'),
        meta: {
            auth: 'sflow.edit',
        }
    },
    {
        name: "sflow_design",
        path: "/sflow/:type/:id/design",
        component: () => import('@/pages/sflow/design/index.vue'),
        meta: {
            auth: 'sflow.design',
        }
    }
]