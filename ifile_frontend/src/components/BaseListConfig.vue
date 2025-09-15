<script setup>
import { ref,watchEffect } from 'vue';
import { post } from '../api.js';
const props = defineProps({
  paths: {
    type: Array,
    default: () => []
  }
});
const emit = defineEmits(['update-paths']);

const fileList = ref([]);

watchEffect(async () => {
  fileList.value = props.paths;
})

const handleRemove = async(tmp) => {
    fileList.value = fileList.value.filter(item => item.id !== tmp.id);
    console.log("remove_dir",tmp.name);
    const res = await post("/api/admin/remove_dir?id=" + tmp.id);
    console.log("remove item res",res);
    emit('update-paths');
};

</script>

<template>
    <div class="bg-orange-50 text-gray-800 mb-6 p-4">
        <div id="homelist" class="w-full max-w-7xl mx-auto">
            <div class="flex items-center justify-between text-xl mb-4 px-2 md:px-0">
                <h1 class="font-bold text-gray-800">目录列表管理</h1>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
                <div v-for="item in fileList" :key="item.id"
                    class="flex flex-col bg-white border-2 border-orange-50 rounded-xl shadow-lg overflow-hidden transition-all duration-300 hover:shadow-xl hover:border-2 hover:border-orange-500"
                    :class="{ 'opacity-50 hover:opacity-100': item.not_exist }">
                    <div class="p-4 flex-grow">
                        <h3 class="text-lg font-semibold text-gray-800">{{ item.name }}</h3>
                    </div>
                    <div class="px-4 py-3 bg-gray-50 border-t border-gray-100 flex justify-end space-x-4">
                        <a v-if="!item.not_exist"
                            :href="`/ui/path.html?base=${item.id}&view=/`" target="_blank"
                            class="px-4 py-2 text-sm font-medium text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors bg-orange-500 hover:bg-orange-600 focus:ring-orange-500"
                        >
                          有效
                        </a>
                        <a v-else title="目录不存在"
                          :href="`/ui/path.html?base=${item.id}&view=/`" target="_blank"
                          class="px-4 py-2 text-sm cursor-pointer font-medium text-gray-500 bg-gray-200 rounded-lg"
                        >
                          无效
                        </a>
                        <a title="移除本目录" @click="handleRemove(item)"
                          class="cursor-pointer px-4 py-2 text-sm font-medium text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors bg-orange-500 hover:bg-orange-600 focus:ring-orange-500"
                        >
                          移除
                        </a>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
