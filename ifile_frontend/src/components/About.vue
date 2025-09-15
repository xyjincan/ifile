<script setup>
import { watchEffect } from "vue";

watchEffect(async () => {
  console.log("about")
})

</script>

<template>
  <div class="bg-gray-100 text-gray-700 min-h-screen flex flex-col">
    <a href="/ui/about" class="text-2xl mt-2 font-bold" >iFIlE About</a>
    <p class="text-lg">ifile tools</p>
    <p class="text-lg">ifile help</p>
    <p class="text-lg">ifile about</p>
  </div>
</template>

<style scoped>
</style>
