import { MessageApi, DialogApi, NotificationApi } from "naive-ui";

import { Router } from 'vue-router'
declare global {
    interface Window {
        router: Router;
        message: MessageApi;
        dialog: DialogApi;
        notification: NotificationApi;
    }
}
