<script setup>
import { ref, onMounted, watch, watchEffect } from "vue";
import FileItem from "./FileItem.vue";
import { post, get, formatBytes } from '../api.js'

const base = ref("");
const view = ref("");
const paths = ref(null);
const isWindows = ref(false);

const uplink = ref({
  name: "上一级",
  base: "",
  view: "",
  isDir: true,
  hidden: false,
  not_delete: true,
});
const total = ref(0);

watchEffect(async () => {
  const ua = navigator.userAgent.toLowerCase();
  if (ua.includes("windows")) {
    //isWindows.value = true;
  }
})

onMounted(async () => {
  //let paras = new URLSearchParams(location.search);// 构造请求参数
  //base.value = paras.get("base")// 根路径ID
  //view.value = paras.get("view")
  //const res = await get("/fs/path?base=" + base.value + "&view=" + view.value);
  // 复用参数
  const res = await get("/fs/path" + location.search);
  //console.log(res);
  base.value = res.base;
  view.value = res.view;
  if (res.code != 200) {
    uplink.value.hidden = true;
    return
  }
  res?.list?.sort((a, b) => {
    return (b.isDir - a.isDir) || a.name.localeCompare(b.name);
  });
  total.value = 0;
  if (res.base) {
    uplink.value.base = res.base;
  }
  if (res.view == "/") {
    uplink.value.hidden = true;
  }
  uplink.value.view = parentPath(res.view);
  if (!res.list) {
    uplink.value.hidden = true;
    return;
  }
  let del_files = new Set();
  const files = res.list.filter(x => {
    const keep = !x.name.endsWith(".ifile_delete");
    if (!keep) {
      del_files.add(x.name.substr(0,x.name.length-".ifile_delete".length));
    }
    return keep;
  });
  for (let x of files) {
    if (del_files.has(x.name) || x.name.endsWith(".ifile_delete")) {
      x.defer_delete = true;
    }
    total.value += x.size;
    x.base = base.value;
    x.path = res.view;
    if (x.isDir) {
      x.view = res.view + x.name + "/";
    } else {
      x.view = res.view + x.name;
    }
    //console.log(x);
  }
  paths.value = files;
})

function clickMark(file) {
  if(!file.click){
    file.click=0;
  }
  file.click += 1;
}

async function handleDelete(file) {
  var r = confirm("确定删除：" + file.name + "?");
  if (!r) {
    return;
  }
  file.defer_delete = true;
  //console.log("handleDelete", file);
  let url = ""
  if (location.port == "8000") {
    url = "http://" + location.hostname + ":8080"
  }
  let api = "/fs/delete";
  let data = {
    "base": file.base,
    "file": file.view,
  }
  let res = await get(api, data);
  //console.log(res);
  //location.reload();
}

function createFileUrl(file) {
  if (!file || !file.base || !file.view) {
    console.error("Invalid file object provided.");
    return null;
  }
  // 1. 对 base 和 view 进行编码
  const encodedBase = encodeURIComponent(file.base);
  let encodedView = encodeURIComponent(file.view);
  //encodedView = encodeURIComponent(encodedView);
  // 2. 构造 UI 路径和查询参数
  //    这部分是固定的，所以可以直接写死
  const finalURL = `/fs/file?base=${encodedBase}&view=/${encodedView}`;
  //console.log(finalURL);
  const baseURL = new URL(window.location.href);
  if (baseURL.port === "8000") {
    baseURL.port = "8080";
  }
  const serverBaseURL = `${baseURL.protocol}//${baseURL.host}`;
  return serverBaseURL + finalURL;
}
// dav url todo drop
function handleFileURL(file) {
  //console.log("file", file);
  //console.log("handleFileURL for file:", file.name);
  // 1. 使用 URL 对象构建基础URL，更健壮且易于维护;URL构造函数会自动处理协议、主机名等
  const baseURL = new URL(window.location.href);
  if (baseURL.port === "8000") {
    baseURL.port = "8080";
  }
  // 移除原始URL的路径、搜索参数和哈希值，只保留协议、域名和端口
  const serverBaseURL = `${baseURL.protocol}//${baseURL.host}`;
  //console.log("Target server base URL:", serverBaseURL);
  const filePath = file.view.startsWith('/') ? file.view : `/${file.view}`;
  // 4. 对路径的各个动态部分分别使用 encodeURIComponent 进行编码
  // 这是处理 '#' 和其他特殊字符的关键步骤！
  const encodedBase = encodeURIComponent(file.base);
  const encodedPath = filePath.split('/').map(segment => encodeURIComponent(segment)).join('/');
  // 注意：上面的 encodedPath 会把路径中的 '/' 也编码掉，这通常不是我们想要的。
  // 更正确的做法是，只对包含动态内容（可能含有特殊字符）的路径部分进行编码。
  // 假设 file.base 和 file.view 的内容是动态且可能包含特殊字符的。
  // 修正后的正确编码方式：
  const safeBase = file.base.split('/').map(segment => encodeURIComponent(segment)).join('/');
  const safeView = file.view.split('/').map(segment => encodeURIComponent(segment)).join('/');
  // 5. 使用模板字符串拼接最终的URL
  // 确保 base 和 view 之间只有一个斜杠
  let finalURL = `${serverBaseURL}/${safeBase}/${safeView}`;
  // 考虑到 file.view 可能已经包含了开头的'/'，我们做一个更智能的拼接
  const combinedPath = [file.base, file.view.startsWith('/') ? file.view.substring(1) : file.view]
    .map(part =>
      part.split('/').map(segment => encodeURIComponent(segment)).join('/')
    )
    .join('/');
  finalURL = `${serverBaseURL}/${combinedPath}`;
  // console.log("Constructed final URL:", finalURL);
  return finalURL;
}

/**
 * 处理文件查看逻辑，构建并打开一个经过正确编码的URL。
 * @param {object} file - 文件对象，期望包含 name, view, 和 base 属性。
 */
function handleView(file) {
  clickMark(file)
  console.log(file);
  console.log("handleView for file:", file.name);
  let finalURL = handleFileURL(file);
  console.log("Constructed final URL:", finalURL);
  window.open(finalURL, "_self");//"_self" "_blank"
}

function handlePlay(event, file) {
  clickMark(file)
  console.log("play_event", file);
  //console.log("play_event",event);
  event.preventDefault();
  //let finalURL = handleFileURL(file);
  let finalURL = createFileUrl(file);
  console.log("vlc final URL:", finalURL);
  navigator?.clipboard?.writeText(finalURL).then(() => {
  }).catch(err => {
    console.error("复制失败: ", err);
  });
  const isAndroid = /Android/i.test(navigator.userAgent);
  if (isAndroid) {
    playInVLC_Android_V2(finalURL, getMimeType(file.name));
  }
  const isIOS = /iPhone|iPad|OS X/i.test(navigator.userAgent);
  if (isIOS) {
    try {
      playInVLC_iOS_3(finalURL);
    } catch (e) {
      alert(e)
    }
  }
}

function playInVLC_iOS_3(videoUrl) {
  const encodedUrl = encodeURIComponent(videoUrl);
  const vlcUrl = `vlc-x-callback://x-callback-url/stream?url=${encodedUrl}`;
  window.location.href = vlcUrl;
}

function playInVLC_Android_V2(videoUrl, mime) {
  const vlcPackage = 'org.videolan.vlc';
  const playStoreUrl = `https://play.google.com/store/apps/details?id=${vlcPackage}`;
  const encodedPlayStoreUrl = encodeURIComponent(playStoreUrl);
  // 直接将完整的 videoUrl 放入 intent: 后面,不再手动拆分 URL。
  // 这让 Android 系统和 VLC 自己去解析完整的 URL，更为稳妥。
  //const mimeType = mime;//
  const mimeType = "video/*";
  let intentUrl = `intent:${videoUrl}#Intent;action=android.intent.action.VIEW;type=${mimeType};package=${vlcPackage};S.browser_fallback_url=${encodedPlayStoreUrl};end`;
  // 打印出来检查一下，确保 intentUrl 看起来是正确的
  //alert("Generated Intent URL: " + intentUrl);
  window.location.href = intentUrl;
}

function getMimeType(fileName) {
  // 提取文件的扩展名
  const extension = fileName.split('.').pop().toLowerCase();
  // 定义各类扩展名对应的类型
  const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg'];
  const videoExtensions = ['mp4', 'mkv', 'avi', 'mov', 'wmv', 'flv', 'webm'];
  const audioExtensions = ['mp3', 'wav', 'aac', 'flac', 'ogg', 'm4a'];
  if (imageExtensions.includes(extension)) {
    return 'image/*';
  } else if (videoExtensions.includes(extension)) {
    return 'video/*';
  } else if (audioExtensions.includes(extension)) {
    return 'audio/*';
  } else {
    return 'file/*';
  }
}

function parentPath(p) {
  // 如果路径以'/'结尾且长度大于1，先移除末尾的'/'
  let path = p.length > 1 && p.endsWith('/') ? p.slice(0, -1) : p;
  // 找到最后一个'/'的位置
  const lastSlashIndex = path.lastIndexOf('/');
  // 如果没有'/'或者'/'就在开头，说明上一级是根目录'/'
  if (lastSlashIndex <= 0) {
    return '/';
  }
  // 截取到最后一个'/'之前的部分
  return path.substring(0, lastSlashIndex);
}





</script>

<template>

  <div class="py-4">
    <a href="/ui/">
      <h1 class="text-md font-bold mb-4">文件列表 {{ formatBytes(total) }}</h1>
    </a>
    <div class="max-w-screen-xl min-w-[18rem]">
      <div v-if="paths" class="flex flex-col xl:grid xl:grid-cols-2 gap-x-2 gap-y-2">
        <FileItem v-if="!uplink.hidden" :file="uplink" @view="handleView" />
        <FileItem v-for="file in paths" :key="file.id" :is_win='isWindows' :file="file" :url=handleFileURL(file)
          @play="handlePlay" @view="handleView" @delete="handleDelete" />
      </div>
      <div v-if="paths?.length == 0" class="text-gray-500 mt-2">
        暂无文件
      </div>
    </div>
  </div>

</template>

<style scoped></style>
