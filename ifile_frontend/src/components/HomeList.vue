<script setup>
import { ref, onMounted, watchEffect } from 'vue';
import { get } from '../api.js';

import BaseListNew from "./BaseListNew.vue";
import BaseListConfig from "./BaseListConfig.vue";

const paths = ref([]);
const showAdd = ref(false);// false true
const showConfigList = ref(false);// false true
const showLocal = ref(false);

watchEffect(async () => {
  if (location.hostname == "localhost" || location.hostname == "127.0.0.1") {
    showLocal.value = true;
  }
})

const loadBaseList = async () => {
  try {
    const res = await get('/fs/paths');
    paths.value = res.list;
  } catch (error) {
    console.error('加载列表失败:', error);
  }
};

const handleDL = (file) => {
  console.log(file);
  let url = location.port === "8000" ? `http://${location.hostname}:8080` : '';
  window.open(`${url}/${file.id}`, "_blank");
};

//
const showPaths = () => {
  showConfigList.value = false;
  showAdd.value = false;
}
const addItem = () => {
  showAdd.value = !showAdd.value;
  if (showAdd.value) {
    showConfigList.value = false;
  }
}
const configItems = () => {
  showConfigList.value = !showConfigList.value;
  if (showConfigList.value) {
    showAdd.value = false;
  }
}

async function handleUpdatePaths() {
  await loadBaseList();
  //console.log('刷新 paths');
}

onMounted(async () => {
  await loadBaseList();
});
</script>
<template>

  <div id="homelist">

    <div class="flex items-center gap-x-8 text-md mb-2">
      <h1 class="font-bold" @click="showPaths" >列表</h1>
      <button v-if="showLocal" @click="addItem"
        class="bg-green-500 cursor-pointer hover:bg-green-600 whitespace-nowrap text-white font-bold py-1 px-2 rounded-md transition-colors duration-200">
        目录配置
      </button>
      <button v-if="showLocal" @click="configItems"
        class="bg-green-500 cursor-pointer hover:bg-green-600 whitespace-nowrap text-white font-bold py-1 px-2 rounded-md transition-colors duration-200">
        目录管理
      </button>
    </div>

    <transition name="fade">
      <BaseListNew v-show="showAdd" :paths='paths' @update-paths="handleUpdatePaths" />
    </transition>
    <transition name="fade">
      <BaseListConfig v-show="showConfigList" :paths='paths' @update-paths="handleUpdatePaths" />
    </transition>

    <div class="max-w-screen-lg 2xl:max-w-screen-xl">
      <div class="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 gap-2">

        <template v-for="file in paths" :key="file.id">
          <!--  -->
          <div v-if="!file.not_exist"
            class="flex flex-col bg-white rounded-xl shadow-lg overflow-hidden transition-all duration-300 hover:shadow-xl border-2 border-orange-50 hover:border-2 hover:border-orange-500">
            <div class="p-4 flex-grow">
              <h3 class="text-lg font-semibold text-gray-800">{{ file.name }}</h3>
            </div>
            <div class="px-4 py-3 bg-gray-50 flex justify-end space-x-4">
              <a :href='"/ui/path.html?base=" + file.id + "&view=/"'
                class="px-3 md:px-6 py-1 text-sm whitespace-nowrap font-medium rounded-lg focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors"
                :class="{
                  'text-white bg-orange-500 hover:bg-orange-600 focus:ring-orange-500': !file.not_exist,
                  'text-gray-400 bg-gray-200 hover:bg-gray-300 hover:text-gray-600': file.not_exist
                }">
                查看
              </a>
              <button @click="handleDL(file)"
                class="px-3 md:px-6 py-1 text-sm font-medium cursor-pointer text-gray-600 rounded-lg focus:outline-none transition-colors bg-gray-200 hover:bg-gray-300 focus:ring-gray-300">
                DAV
              </button>
            </div>
          </div>

        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-to,
.fade-leave-from {
  opacity: 1;
}
</style>
