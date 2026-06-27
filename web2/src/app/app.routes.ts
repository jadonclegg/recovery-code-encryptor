import { Route, Routes } from '@angular/router';
import { EncryptionPage } from './encryption-page/encryption-page';

export const routes: Routes = [
    <Route>{
        pathMatch: "prefix",
        path: "",
        component: EncryptionPage
    }
];
