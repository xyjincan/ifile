import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'

import './style.css'

import App from './App.vue'
// VUE路由组件
import HomeList from './components/HomeList.vue'
import About from './components/About.vue'

const routes = [
    { path: '/', name: 'Home', component: HomeList },
    { path: '/about', name: 'About', component: About }
]

// createWebHashHistory createWebHistory("/ui")
// 后端接口优化：/ui/*
const router = createRouter({
    history: createWebHistory("/ui"),
    routes
})

router.afterEach(async (to, from) => {
    console.log("index: " + to.path);
});

const app = createApp(App)
app.use(router)
app.mount('#app')
