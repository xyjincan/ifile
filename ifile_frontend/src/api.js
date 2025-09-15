let API_URL = "";

async function get(url, params) {
  if (url.indexOf("?") >= 0) {
    
  } else {
    const searchParams = new URLSearchParams();
    if (params != undefined) {
      for (let x in params) {
        searchParams.set(x, params[x])
      }
    }
    if(searchParams.size != 0){
      url = url + '?' + searchParams.toString();
    }
  }
  url = API_URL + url;
  let res = await (await fetch(url)).json();
  return res;
}

async function post(url,params){
  url = API_URL + url;
  let res = {"code":-1}
  try {
    res = await (await fetch(url,{
      headers: {
          'Accept': 'application/json, text/plain, */*',
          'Content-Type': 'application/json'
      },
      method: "POST",
      body: JSON.stringify(params)
  })).json();
  } catch (error) {
    console.error(error);
  }
  return res;
}

function formatBytes(bytes, decimals = 2) {
  // console.log("todo:",bytes)
  if (bytes == 0) return '0 Bytes';
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


function generateNextID(existingSet) {
  const charSet = Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i));
  // 先检查所有单字符
  for (let i = 0; i < charSet.length; i++) {
    const id = charSet[i];
    if (!existingSet.has(id)) {
      return id;
    }
  }
  // 单字符都用完，检查所有两字符组合
  for (let i = 0; i < charSet.length; i++) {
    for (let j = 0; j < charSet.length; j++) {
      const id = charSet[i] + charSet[j];
      if (!existingSet.has(id)) {
        return id;
      }
    }
  }
  // 如果所有组合都用完
  return null;
}

/**
 * 格式化文件大小，并确保小数点前的数字不超过两位。
 */
function formatBytesWithLimit(bytes, decimals = 2) {
  if (!+bytes) return '0 Bytes';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  // 2. 初始化单位索引和用于计算的数值
  let i = 0;
  let formattedValue = bytes;
  while (formattedValue >= 100 && i < sizes.length - 1) {
    formattedValue /= k; // 将数值除以1024
    i++;                 // 单位索引增加，例如从 KB -> MB
  }
  return `${parseFloat(formattedValue.toFixed(dm))} ${sizes[i]}`;
}

export { get, post, formatBytes,generateNextID };

