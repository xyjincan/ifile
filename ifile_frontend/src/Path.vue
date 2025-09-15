<script setup>
import { ref, onMounted, watch, watchEffect } from "vue";

import FileList from "./components/FileList.vue";

watchEffect(async () => {
})

</script>

<template>

  <div class="bg-gray-100 text-gray-700 min-h-screen flex flex-col p-1 px-[6px]">
    <a href="/ui/" class="text-2xl mt-2 font-bold">iFIlE</a>
    <p class="text-lg text-cyan-400 ">local web files</p>

    <FileList />

  </div>

</template>

<style scoped></style>
