const fs = require('fs');
const path = require('path');

const sqlPath = path.join(__dirname, 'myself3dblog_2026-10-05_01-13-05_mysql_data_vw18p.sql');
const content = fs.readFileSync(sqlPath, 'utf8');

const articleInsertsMatch = content.match(/INSERT INTO `articles` VALUES\s*([\s\S]+?);(?=\s*(?:INSERT|UNLOCK|DROP|CREATE|ALTER|$))/);
const rawValues = articleInsertsMatch[1];

function parseSqlRows(str) {
  const rows = [];
  let inRow = false;
  let inString = false;
  let escape = false;
  let currentVal = '';
  let currentRow = [];

  for (let i = 0; i < str.length; i++) {
    const char = str[i];
    if (!inRow) {
      if (char === '(') {
        inRow = true;
        currentRow = [];
        currentVal = '';
        inString = false;
      }
    } else {
      if (inString) {
        if (escape) {
          currentVal += char;
          escape = false;
        } else if (char === '\\') {
          escape = true;
        } else if (char === "'") {
          if (str[i + 1] === "'") {
            currentVal += "'";
            i++;
          } else {
            inString = false;
          }
        } else {
          currentVal += char;
        }
      } else {
        if (char === "'") {
          inString = true;
        } else if (char === ',') {
          currentRow.push(currentVal.trim());
          currentVal = '';
        } else if (char === ')') {
          currentRow.push(currentVal.trim());
          rows.push(currentRow);
          currentRow = [];
          inRow = false;
        } else {
          currentVal += char;
        }
      }
    }
  }
  return rows;
}

const rows = parseSqlRows(rawValues);

function extractLinks(text) {
  const links = [];
  const regex = /https?:\/\/[^\s"'<>\\]+/g;
  let m;
  while ((m = regex.exec(text)) !== null) {
    links.push(m[0]);
  }
  return links;
}

function parseResourceDownloads(resContent) {
  const downloads = [];
  if (!resContent) return downloads;

  // 1. Google Drive
  const gMatch = resContent.match(/https?:\/\/(?:drive\.google\.com)[^\s"'<>\\]+/);
  if (gMatch) {
    downloads.push({
      platform: 'google',
      url: gMatch[0],
      extract_code: ''
    });
  }

  // 2. Quark
  const qMatch = resContent.match(/https?:\/\/(?:pan\.quark\.cn)[^\s"'<>\\]+/);
  if (qMatch) {
    let qCode = '';
    // Look for pwd=... or 提取码: XXXX
    const pwdParam = qMatch[0].match(/pwd=([a-zA-Z0-9]+)/);
    const codeMatch = resContent.match(/(?:提取码|密码)[^\w]*([a-zA-Z0-9]{4,8})/);
    if (pwdParam) qCode = pwdParam[1];
    else if (codeMatch) qCode = codeMatch[1];

    downloads.push({
      platform: 'quark',
      url: qMatch[0],
      extract_code: qCode
    });
  }

  // 3. Baidu
  const bMatch = resContent.match(/https?:\/\/(?:pan\.baidu\.com)[^\s"'<>\\]+/);
  if (bMatch) {
    let bCode = '';
    const pwdParam = bMatch[0].match(/pwd=([a-zA-Z0-9]+)/);
    const codeMatch = resContent.match(/(?:提取码|密码)[^\w]*([a-zA-Z0-9]{4,8})/);
    if (pwdParam) bCode = pwdParam[1];
    else if (codeMatch) bCode = codeMatch[1];

    downloads.push({
      platform: 'baidu',
      url: bMatch[0],
      extract_code: bCode
    });
  }

  return downloads;
}

function parseVideos(contentStr) {
  const videos = [];
  const urls = extractLinks(contentStr);
  const seen = new Set();

  urls.forEach(url => {
    if (seen.has(url)) return;
    const lower = url.toLowerCase();
    if (lower.includes('tiktok.com')) {
      seen.add(url);
      videos.push({ platform: 'tiktok', url, title: 'TikTok 演示视频' });
    } else if (lower.includes('douyin.com')) {
      seen.add(url);
      videos.push({ platform: 'douyin', url, title: '抖音 演示视频' });
    } else if (lower.includes('xiaohongshu.com') || lower.includes('xhslink.com')) {
      seen.add(url);
      videos.push({ platform: 'xiaohongshu', url, title: '小红书 实测笔记' });
    } else if (lower.includes('bilibili.com')) {
      seen.add(url);
      videos.push({ platform: 'bilibili', url, title: 'B站 高清演示' });
    } else if (lower.includes('youtube.com') || lower.includes('youtu.be')) {
      seen.add(url);
      videos.push({ platform: 'youtube', url, title: 'YouTube 演示' });
    }
  });

  return videos;
}

const parsedProducts = rows.map((r, i) => {
  const id = r[0];
  const title_zh = r[1] || '精美壁纸';
  const title_en = r[2] || '';
  const category = r[3] || '苹果壁纸';
  const contentStr = r[6] || '';
  const cover_image = r[8] || '';
  const price = parseFloat(r[11]) || parseFloat(r[23]) || 0;
  const views = parseInt(r[13]) || 0;
  const created_at = r[14] || new Date().toISOString();
  const resource_content = r[16] || '';

  const videos = parseVideos(contentStr);
  const downloads = parseResourceDownloads(resource_content);

  let primaryPlatform = 'douyin';
  let primaryVideo = '';
  if (videos.length > 0) {
    primaryPlatform = videos[0].platform;
    primaryVideo = videos[0].url;
  }

  return {
    legacy_id: id,
    title_zh,
    title_en,
    category,
    price,
    verify_code: r[15] || '',
    cover_image,
    preview_platform: primaryPlatform,
    preview_video_url: primaryVideo,
    preview_videos: videos,
    downloads,
    description: contentStr,
    views_count: views,
    created_at
  };
});

fs.writeFileSync(path.join(__dirname, 'legacy_products.json'), JSON.stringify(parsedProducts, null, 2), 'utf8');
console.log(`Saved ${parsedProducts.length} converted products to legacy_products.json!`);

let withVideos = parsedProducts.filter(p => p.preview_videos.length > 0).length;
let withDownloads = parsedProducts.filter(p => p.downloads.length > 0).length;
console.log(`Products with videos: ${withVideos}`);
console.log(`Products with downloads: ${withDownloads}`);
