<script setup>
import { watchEffect } from "vue";

watchEffect(async () => {
})

</script>

<template>

  <div class="bg-orange-50 text-gray-800 min-h-screen flex flex-col p-1 px-[0.5rem] overflow-x-auto">
    <a href="/ui/" class="text-2xl mt-2 font-bold text-orange-600" >iFIlE</a>
    <p title="请勿映射到公网IP 端口；不要在公共网络中使用；" class="text-lg text-cyan-400 text-gray-600 ">本地网络文件目录</p>

    <RouterView />

  </div>

</template>

<style scoped>
</style>
