
<script setup>
import { ref } from 'vue';
import { post,generateNextID } from '../api.js';

const props = defineProps({
  paths: {
    type: Array,
    default: () => []
  }
});
const emit = defineEmits(['update-paths']);

// 下一个条目的唯一 ID，用于给 v-for 提供稳定的 :key
let nextId = 1;
const items = ref([
    {
        "id": 0,
        "path": "",
        "name": ""
    }
]);

async function saveItems() {
    const validItems = items.value.filter(item => item.path.trim() !== "");
    const idSet = new Set(["ui", "fs"]);
    props.paths.forEach(item => idSet.add(item.id));
    validItems.forEach(item => {
        let nid = generateNextID(idSet);
        item.id = nid;
        idSet.add(nid);
    } );
    for (const item of validItems) {
      console.log("add item",item);
      const res = await post("/api/admin/add_dir",item);
      console.log("add item res",res);
    }
    emit('update-paths');
}

// 添加一个新条目的方法
const addItem = () => {
    items.value.push({
        id: nextId++,
        path: '',
        name: ''
    });
};

// 根据索引删除一个条目的方法
const removeItem = (index) => {
    items.value.splice(index, 1);
    // 
};
</script>
<template>
    <div class="px-4 py-4 max-w-2xl mx-auto bg-white rounded-xl shadow-md space-y-4 mb-4">
        <div class="flex justify-between items-center mb-2">
            <h2 class="text-xl font-bold text-gray-900">新增列表</h2>

            <div class="flex gap-x-4">
                <button @click="addItem"
                    class="bg-orange-500 hover:bg-orange-600 cursor-pointer text-white font-bold py-1 px-2 rounded-lg transition-colors duration-200">
                    增加条目
                </button>
                <button @click="saveItems"
                    class="bg-orange-500 hover:bg-orange-600 cursor-pointer text-white font-bold py-1 px-2 rounded-lg transition-colors duration-200">
                    保存配置
                </button>
            </div>

        </div>

        <div v-if="items.length > 0" class="space-y-4">
            <div v-for="(item, index) in items" :key="item.id"
                class="px-4 py-2 rounded-lg flex flex-col sm:flex-row sm:items-center sm:gap-x-4 space-y-2 sm:space-y-0">

                <div class="flex-grow space-y-4">
                    <div class="flex flex-col">
                        <label :for="`path-${item.id}`" class="text-sm font-medium text-gray-700 mb-1">本地路径:</label>
                        <input :id="`path-${item.id}`" v-model="item.path" type="text" autocomplete="off" placeholder="请输入本地路径"
                            class="block w-full rounded-md border border-gray-300 bg-white px-3 py-2 shadow-sm focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition duration-150 ease-in-out" />
                    </div>

                    <div class="flex flex-col">
                        <label :for="`name-${item.id}`" class="text-sm font-medium text-gray-700 mb-1">名称 (可选):</label>
                        <input :id="`name-${item.id}`" v-model="item.name" type="text" placeholder="请输入名称"
                            class="block w-full rounded-md border border-gray-300 bg-white px-3 py-2 shadow-sm focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition duration-150 ease-in-out" />
                    </div>
                </div>

                <button @click="removeItem(index)"
                    class="bg-orange-500 hover:bg-orange-600 max-w-[6rem] cursor-pointer text-white font-bold py-2 px-4 rounded-lg transition-colors duration-200">
                    删除
                </button>
            </div>
        </div>
        <p v-else class="text-gray-500 text-center italic">
            列表为空，请点击“增加条目”按钮。
        </p>

        <div v-if="false" class="mt-6 p-4 bg-gray-50 rounded-lg">
            <h3 class="text-md font-semibold text-gray-800">当前数据:</h3>
            <pre class="whitespace-pre-wrap text-sm text-gray-600">{{ JSON.stringify(items, null, 2) }}</pre>
        </div>
    </div>
</template>
<style scoped>
/* 如果你需要额外的自定义样式，可以在这里添加 */
</style>