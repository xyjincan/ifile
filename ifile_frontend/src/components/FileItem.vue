<script setup>
import { ref, onMounted } from 'vue'

const prop = defineProps({
  is_win: Boolean,
  file: Object,
  url: String,
  required: true,
})
const emit = defineEmits(['play', 'view', 'delete'])

function handlePlay(event, file) {
  emit('play', event, file);
}

onMounted(async () => {
})

function formatGolangTime(str) {
  const t = new Date(str);
  if (isNaN(t.getTime())) {
    return "";
  }
  const formattedDate = t.toISOString().slice(0, 10);
  return formattedDate;
}

function formatBytes(bytes, decimals = 2) {
  //console.log("todo bytes",bytes);
  if (!bytes) {
    return '';//0 Bytes
  }
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  // 如果结果是MB或GB，则显示相应的单位
  if (i === 2) { // MB
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' MB';
  } else if (i === 3) { // GB
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' GB';
  } else { // 其他单位，如Bytes, KB, TB等
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
  }
}

function createUIUrl(file) {
  // 确保 file 对象和其属性存在
  if (!file || !file.base || !file.view) {
    console.error("Invalid file object provided.");
    return null;
  }
  const encodedBase = encodeURIComponent(file.base);
  let encodedView = encodeURIComponent(file.view);
  const finalURL = `/ui/path.html?base=${encodedBase}&view=/${encodedView}&dir=${file.isDir}`;
  return finalURL;
}

function createFileUrl(file) {
  if (!file || !file.base || !file.view) {
    console.error("Invalid file object provided.");
    return null;
  }
  const encodedBase = encodeURIComponent(file.base);
  let encodedView = encodeURIComponent(file.view);
  const finalURL = `/fs/file?base=${encodedBase}&view=/${encodedView}&dir=${file.isDir}`;
  return finalURL;
}

function clickMark(file) {
  if(!file.click){
    file.click=0;
  }
  file.click += 1;
}

</script>
<template>

  <div class="max-w-screen-md w-full flex items-center justify-between px-2 py-1 rounded shadow-sm">
    <a :title="file?.name" :href='file.isDir ? createUIUrl(file) : createFileUrl(file)'
      @click="clickMark(file)"
      class="flex-1 min-w-0 flex flex-col">
      <p class="min-w-[3rem] text-sm break-words font-medium w-full" 
        :class='{ "text-yellow-500": file.defer_delete }'
      >
        {{ file?.name }}
      </p>
      <span class="flex items-center text-xs text-gray-400 select-none">
        <span v-if="file?.size" class="pr-2" >{{ formatBytes(file?.size) }}</span>
        <span class="pr-2">{{ formatGolangTime(file?.time) }}</span>
        <span v-if="file.defer_delete" 
          title="延时删除"
          class="text-xs text-yellow-500 select-none">
          X
        </span>
        <span v-else-if="file.click" title="like"
          class="text-gray-400 select-none"
            :class='{ "text-red-400": file.click >= 3 }'
           >
          <svg width="10" height="10" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-heart-icon lucide-heart"><path d="M2 9.5a5.5 5.5 0 0 1 9.591-3.676.56.56 0 0 0 .818 0A5.49 5.49 0 0 1 22 9.5c0 2.29-1.5 4-3 5.5l-5.492 5.313a2 2 0 0 1-3 .019L5 15c-1.5-1.5-3-3.2-3-5.5"/>
          </svg>
        </span>
      </span>
    </a>
    <div class="ml-2 flex space-x-2 select-none">
      <!-- 目录 -->
      <a v-if="file?.isDir" title="enter" :href='createUIUrl(file)'
        class="px-1 py-1 bg-yellow-400 text-white rounded hover:bg-blue-600">
        <svg class="size-6" viewBox="0 0 32 32" title="Open Folder">
          <g data-name="Layer 39" id="Layer_39">
            <path class="cls-1"
              d="M4,28a3,3,0,0,1-3-3V7A3,3,0,0,1,4,4h7a1,1,0,0,1,.77.36L14.8,8H27a1,1,0,0,1,0,2H14.33a1,1,0,0,1-.76-.36L10.53,6H4A1,1,0,0,0,3,7V25a1,1,0,0,0,1,1,1,1,0,0,1,0,2Z" />
            <path class="cls-1"
              d="M25.38,28H4a1,1,0,0,1-1-1.21l3-14A1,1,0,0,1,7,12H30a1,1,0,0,1,1,1.21L28.32,25.63A3,3,0,0,1,25.38,28ZM5.24,26H25.38a1,1,0,0,0,1-.79L28.76,14h-21Z" />
          </g>
        </svg>
      </a>
      <!-- 文件 -->
      <!-- 类型 处理方式 图标选择 -->
      <a v-else-if="!is_win" title="open with" :href='createFileUrl(file)' @click="handlePlay($event, file)"
        class="place-self-center px-1 py-1 bg-gray-50 text-white rounded hover:bg-orange-100">
        <svg class="size-6" viewBox="0 0 32 32">
          <path fill="#e67600"
            d="M16 32c-8.836 0-16-7.164-16-16s7.164-16 16-16 16 7.163 16 16c0 8.836-7.164 16-16 16zM16 4c-6.627 0-12 5.373-12 12s5.373 12 12 12 12-5.373 12-12-5.373-12-12-12zM12 10l10 6-10 6v-12z">
          </path>
        </svg>
      </a>
      <a v-if="!file.isDir" title="open file" :href='createFileUrl(file)'
        @click="clickMark(file)"
        class="px-1 py-1 bg-gray-50 text-white rounded hover:bg-blue-100">
        <svg class="size-6" viewBox="0 0 20 20">
          <path fill="#7DD3FC"
            d="M0 4c0-1.1.9-2 2-2h16a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V4zm6 0v12h8V4H6zM2 5v2h2V5H2zm0 4v2h2V9H2zm0 4v2h2v-2H2zm14-8v2h2V5h-2zm0 4v2h2V9h-2zm0 4v2h2v-2h-2zM8 7l5 3-5 3V7z" />
        </svg>
      </a>
      <button title="webdav" v-if="file.isDir"
        class="hover:cursor-pointer rounded place-self-center px-1 py-1 bg-amber-400 hover:bg-amber-400"
        @click="$emit('view', file)" >
        <svg class="size-6" viewBox="0 0 20 20">
          <path
            d="M0 4c0-1.1.9-2 2-2h16a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V4zm6 0v12h8V4H6zM2 5v2h2V5H2zm0 4v2h2V9H2zm0 4v2h2v-2H2zm14-8v2h2V5h-2zm0 4v2h2V9h-2zm0 4v2h2v-2h-2zM8 7l5 3-5 3V7z" />
        </svg>
      </button>
      <!-- v-if="file?.size" v-if="file" -->
      <button title="delete_file" v-if="!file?.not_delete"
        class="hover:cursor-pointer px-1 py-1 bg-white text-red-500 hover:text-red-500 rounded "
        @click="$emit('delete', file)" >
        <svg class="size-6" fill="currentColor" version="1.1" viewBox="0 0 40 40">
          <g>
            <path
              d="M28,40H11.8c-3.3,0-5.9-2.7-5.9-5.9V16c0-0.6,0.4-1,1-1s1,0.4,1,1v18.1c0,2.2,1.8,3.9,3.9,3.9H28c2.2,0,3.9-1.8,3.9-3.9V16   c0-0.6,0.4-1,1-1s1,0.4,1,1v18.1C33.9,37.3,31.2,40,28,40z" />
          </g>
          <g>
            <path
              d="M33.3,4.9h-7.6C25.2,2.1,22.8,0,19.9,0s-5.3,2.1-5.8,4.9H6.5c-2.3,0-4.1,1.8-4.1,4.1S4.2,13,6.5,13h26.9   c2.3,0,4.1-1.8,4.1-4.1S35.6,4.9,33.3,4.9z M19.9,2c1.8,0,3.3,1.2,3.7,2.9h-7.5C16.6,3.2,18.1,2,19.9,2z M33.3,11H6.5   c-1.1,0-2.1-0.9-2.1-2.1c0-1.1,0.9-2.1,2.1-2.1h26.9c1.1,0,2.1,0.9,2.1,2.1C35.4,10.1,34.5,11,33.3,11z" />
          </g>
          <g>
            <path d="M12.9,35.1c-0.6,0-1-0.4-1-1V17.4c0-0.6,0.4-1,1-1s1,0.4,1,1v16.7C13.9,34.6,13.4,35.1,12.9,35.1z" />
          </g>
          <g>
            <path d="M26.9,35.1c-0.6,0-1-0.4-1-1V17.4c0-0.6,0.4-1,1-1s1,0.4,1,1v16.7C27.9,34.6,27.4,35.1,26.9,35.1z" />
          </g>
          <g>
            <path d="M19.9,35.1c-0.6,0-1-0.4-1-1V17.4c0-0.6,0.4-1,1-1s1,0.4,1,1v16.7C20.9,34.6,20.4,35.1,19.9,35.1z" />
          </g>
        </svg>
      </button>
    </div>
  </div>
</template>
<style scoped></style>
