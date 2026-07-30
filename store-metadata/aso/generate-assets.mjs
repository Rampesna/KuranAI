import { mkdir, writeFile } from "node:fs/promises";
import { execFile } from "node:child_process";
import { promisify } from "node:util";
import path from "node:path";

const run = promisify(execFile);
const root = path.resolve(import.meta.dirname, "assets");

const locales = {
  "tr-TR": {
    rtl: false,
    eyebrow: "KURANAI",
    slides: [
      {
        title: ["Kur’an’da ne", "söylendiğini sor"],
        subtitle: "Aklına takılan konuyu ayet ve kaynaklarla araştır.",
      },
      {
        title: ["Ayet, meal,", "tefsir ve hadis"],
        subtitle: "Kaynak odağını seç; düzenli ve anlaşılır yanıt al.",
      },
      {
        title: ["Konuş,", "kayıt kendi dursun"],
        subtitle: "Sorunu sesli sor; sessizlik algılanınca kayıt bitsin.",
      },
      {
        title: ["Üç dil,", "tek deneyim"],
        subtitle: "Türkçe, English ve tam RTL destekli العربية.",
      },
    ],
    sourceQuran: "Kur’an",
    sourceTafsir: "Tefsir",
    sourceHadith: "Hadis",
    askPlaceholder: "Kur’an hakkında sor...",
    dailyVerse: "Günün Ayeti",
    verseMeaning: "Şüphesiz zorlukla beraber kolaylık vardır.",
    prayerTitle: "Namaz Vakitleri",
    prayerPlace: "İstanbul",
    question: "Sabır hakkında Kur’an ne söyler?",
    answerTitle: "Kur’an’da sabır",
    answerLines: [
      "Sabır; zorluk karşısında direnç,",
      "dua ve doğru davranışla birlikte",
      "ele alınır.",
    ],
    sources: "İlgili kaynaklar",
    listening: "Dinliyorum…",
    autoStop: "Konuşma bitince otomatik durur",
    voiceHint: "Sorunu doğal biçimde söyle",
    historyTitle: "Sohbet Geçmişi",
    historyRows: [
      "Sabır hakkında Kur’an ne söyler?",
      "Rahman suresindeki mesaj nedir?",
      "Yolculukta namaz nasıl anlatılır?",
    ],
    languageTitle: "Dili istediğin an değiştir",
    featureTitle: ["Kur’an’da ne söylendiğini", "ayet ve kaynaklarla keşfet"],
    featureSubtitle: "Meal, tefsir, hadis, sesli soru ve namaz vakitleri",
  },
  "en-US": {
    rtl: false,
    eyebrow: "KURANAI",
    slides: [
      {
        title: ["Ask what", "the Quran says"],
        subtitle: "Explore your question through verses and relevant sources.",
      },
      {
        title: ["Verses, tafsir", "& hadith"],
        subtitle: "Choose a source focus and get a clear, structured answer.",
      },
      {
        title: ["Speak.", "Recording stops itself."],
        subtitle: "Ask naturally; silence detection ends the recording.",
      },
      {
        title: ["Three languages,", "one experience"],
        subtitle: "Turkish, English and Arabic with complete RTL support.",
      },
    ],
    sourceQuran: "Quran",
    sourceTafsir: "Tafsir",
    sourceHadith: "Hadith",
    askPlaceholder: "Ask about the Quran...",
    dailyVerse: "Verse of the Day",
    verseMeaning: "Surely with hardship comes ease.",
    prayerTitle: "Prayer Times",
    prayerPlace: "Istanbul",
    question: "What does the Quran say about patience?",
    answerTitle: "Patience in the Quran",
    answerLines: [
      "Patience is presented together with",
      "prayer, resilience and doing what is",
      "right through hardship.",
    ],
    sources: "Related sources",
    listening: "Listening…",
    autoStop: "Stops automatically when you finish",
    voiceHint: "Ask your question naturally",
    historyTitle: "Chat History",
    historyRows: [
      "What does the Quran say about patience?",
      "What is the message of Surah Ar-Rahman?",
      "How is prayer described while traveling?",
    ],
    languageTitle: "Switch languages at any time",
    featureTitle: ["Explore what the Quran says", "through verses and sources"],
    featureSubtitle: "Translation, tafsir, hadith, voice questions and prayer times",
  },
  ar: {
    rtl: true,
    eyebrow: "KURANAI",
    slides: [
      {
        title: ["اسأل عمّا", "يقوله القرآن"],
        subtitle: "ابحث في سؤالك من خلال الآيات والمصادر ذات الصلة.",
      },
      {
        title: ["آيات وتفسير", "وحديث"],
        subtitle: "اختر نوع المصدر واحصل على إجابة واضحة ومنظّمة.",
      },
      {
        title: ["تحدّث، وسيتوقف", "التسجيل تلقائياً"],
        subtitle: "اطرح سؤالك طبيعياً؛ ينتهي التسجيل عند اكتشاف الصمت.",
      },
      {
        title: ["ثلاث لغات،", "وتجربة واحدة"],
        subtitle: "العربية والتركية والإنجليزية مع دعم كامل لاتجاه RTL.",
      },
    ],
    sourceQuran: "القرآن",
    sourceTafsir: "التفسير",
    sourceHadith: "الحديث",
    askPlaceholder: "اسأل عن القرآن…",
    dailyVerse: "آية اليوم",
    verseMeaning: "فَإِنَّ مَعَ الْعُسْرِ يُسْرًا",
    prayerTitle: "مواقيت الصلاة",
    prayerPlace: "إسطنبول",
    question: "ماذا يقول القرآن عن الصبر؟",
    answerTitle: "الصبر في القرآن",
    answerLines: [
      "يرتبط الصبر في القرآن بالصلاة",
      "والثبات والعمل الصالح عند",
      "مواجهة الشدائد.",
    ],
    sources: "المصادر ذات الصلة",
    listening: "جارٍ الاستماع…",
    autoStop: "يتوقف تلقائياً عند انتهاء الكلام",
    voiceHint: "اطرح سؤالك بصورة طبيعية",
    historyTitle: "سجل المحادثات",
    historyRows: [
      "ماذا يقول القرآن عن الصبر؟",
      "ما رسالة سورة الرحمن؟",
      "كيف وُصفت الصلاة في السفر؟",
    ],
    languageTitle: "غيّر اللغة في أي وقت",
    featureTitle: ["اكتشف ما يقوله القرآن", "من خلال الآيات والمصادر"],
    featureSubtitle: "المعاني والتفسير والحديث والأسئلة الصوتية ومواقيت الصلاة",
  },
};

const specs = {
  "app-store": { width: 1320, height: 2868 },
  "google-play": { width: 1080, height: 1920 },
};

const escapeXml = (value) =>
  String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");

function textBlock({
  lines,
  x,
  y,
  fontSize,
  lineHeight,
  fill = "#fff8e8",
  weight = 700,
  anchor = "start",
  family = "Arial, sans-serif",
  opacity = 1,
}) {
  return `<text x="${x}" y="${y}" text-anchor="${anchor}" fill="${fill}" fill-opacity="${opacity}" font-family="${family}" font-size="${fontSize}" font-weight="${weight}">
${lines
  .map(
    (line, index) =>
      `<tspan x="${x}" dy="${index === 0 ? 0 : lineHeight}">${escapeXml(line)}</tspan>`,
  )
  .join("\n")}
</text>`;
}

function baseDefs(id) {
  return `<defs>
    <linearGradient id="bg-${id}" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#06130f"/>
      <stop offset="0.52" stop-color="#0d1813"/>
      <stop offset="1" stop-color="#050907"/>
    </linearGradient>
    <radialGradient id="glow-${id}" cx="50%" cy="15%" r="70%">
      <stop offset="0" stop-color="#c99a18" stop-opacity=".22"/>
      <stop offset=".55" stop-color="#c99a18" stop-opacity=".04"/>
      <stop offset="1" stop-color="#c99a18" stop-opacity="0"/>
    </radialGradient>
    <radialGradient id="mic-${id}" cx="50%" cy="50%" r="50%">
      <stop offset="0" stop-color="#ffd94c"/>
      <stop offset=".58" stop-color="#e7b71d"/>
      <stop offset="1" stop-color="#9d7410"/>
    </radialGradient>
    <filter id="shadow-${id}" x="-25%" y="-20%" width="150%" height="160%">
      <feDropShadow dx="0" dy="26" stdDeviation="30" flood-color="#000" flood-opacity=".58"/>
    </filter>
    <filter id="soft-${id}" x="-80%" y="-80%" width="260%" height="260%">
      <feGaussianBlur stdDeviation="28"/>
    </filter>
  </defs>`;
}

function sourceTabs(locale, width, y, compact = false) {
  const labels = [locale.sourceQuran, locale.sourceTafsir, locale.sourceHadith];
  const gap = width * 0.018;
  const tabWidth = (width - gap * 2) / 3;
  const height = compact ? 54 : 64;
  return labels
    .map((label, index) => {
      const x = index * (tabWidth + gap);
      const active = index === 0;
      return `<g transform="translate(${x} ${y})">
        <rect width="${tabWidth}" height="${height}" rx="${height / 2}" fill="${active ? "#f0bd24" : "#211f1b"}" stroke="${active ? "#f0bd24" : "#39362f"}"/>
        <text x="${tabWidth / 2}" y="${height * 0.64}" text-anchor="middle" fill="${active ? "#10120f" : "#ece7da"}" font-family="Arial, sans-serif" font-size="${compact ? 22 : 25}" font-weight="700">${escapeXml(label)}</text>
      </g>`;
    })
    .join("");
}

function micIcon(cx, cy, size, id) {
  return `<g transform="translate(${cx} ${cy})">
    <circle r="${size * 0.64}" fill="#d7a515" opacity=".13" filter="url(#soft-${id})"/>
    <circle r="${size * 0.52}" fill="url(#mic-${id})"/>
    <rect x="${-size * 0.11}" y="${-size * 0.24}" width="${size * 0.22}" height="${size * 0.34}" rx="${size * 0.11}" fill="#07130f"/>
    <path d="M ${-size * 0.2} ${size * 0.02} Q 0 ${size * 0.3} ${size * 0.2} ${size * 0.02}" fill="none" stroke="#07130f" stroke-width="${size * 0.055}" stroke-linecap="round"/>
    <path d="M 0 ${size * 0.27} V ${size * 0.38}" stroke="#07130f" stroke-width="${size * 0.055}" stroke-linecap="round"/>
    <path d="M ${-size * 0.15} ${size * 0.38} H ${size * 0.15}" stroke="#07130f" stroke-width="${size * 0.055}" stroke-linecap="round"/>
  </g>`;
}

function phoneChrome({ x, y, width, height, id, content }) {
  const radius = width * 0.09;
  return `<g filter="url(#shadow-${id})">
    <rect x="${x - 8}" y="${y - 8}" width="${width + 16}" height="${height + 16}" rx="${radius + 8}" fill="#030504" stroke="#b98c18" stroke-opacity=".52" stroke-width="4"/>
    <clipPath id="clip-${id}"><rect x="${x}" y="${y}" width="${width}" height="${height}" rx="${radius}"/></clipPath>
    <g clip-path="url(#clip-${id})">
      <rect x="${x}" y="${y}" width="${width}" height="${height}" fill="#090805"/>
      ${content}
    </g>
  </g>`;
}

function sceneHome(locale, box, id) {
  const { x, y, width: w, height: h } = box;
  const pad = w * 0.065;
  const anchor = locale.rtl ? "end" : "start";
  const tx = locale.rtl ? x + w - pad : x + pad;
  const tabsY = y + h * 0.18;
  const cardY = y + h * 0.31;
  const cardW = w - pad * 2;
  return `
    <text x="${tx}" y="${y + h * 0.075}" text-anchor="${anchor}" fill="#e7b72b" font-family="Arial, sans-serif" font-size="${w * 0.045}" font-weight="800">KuranAI</text>
    <text x="${tx}" y="${y + h * 0.12}" text-anchor="${anchor}" fill="#a9ada8" font-family="Arial, sans-serif" font-size="${w * 0.025}">${escapeXml(locale.askPlaceholder)}</text>
    <g transform="translate(${x + pad} 0)">${sourceTabs(locale, cardW, tabsY, true)}</g>
    <rect x="${x + pad}" y="${cardY}" width="${cardW}" height="${h * 0.27}" rx="${w * 0.045}" fill="#171510" stroke="#3c3525"/>
    <circle cx="${x + w / 2}" cy="${cardY + h * 0.105}" r="${w * 0.075}" fill="#3a3216"/>
    <text x="${x + w / 2}" y="${cardY + h * 0.13}" text-anchor="middle" fill="#f2bd21" font-family="Geeza Pro, Arial, sans-serif" font-size="${w * 0.09}" font-weight="700">ق</text>
    <text x="${x + w / 2}" y="${cardY + h * 0.195}" text-anchor="middle" fill="#fff8e8" font-family="Arial, sans-serif" font-size="${w * 0.033}" font-weight="700">${escapeXml(locale.askPlaceholder)}</text>
    <rect x="${x + pad}" y="${y + h * 0.62}" width="${cardW}" height="${h * 0.19}" rx="${w * 0.04}" fill="#171510" stroke="#3c3525"/>
    <text x="${tx}" y="${y + h * 0.665}" text-anchor="${anchor}" fill="#e7b72b" font-family="Arial, sans-serif" font-size="${w * 0.03}" font-weight="700">${escapeXml(locale.dailyVerse)}</text>
    <text x="${tx}" y="${y + h * 0.72}" text-anchor="${anchor}" fill="#fff8e8" font-family="Geeza Pro, Arial, sans-serif" font-size="${w * 0.037}" font-weight="600">${escapeXml(locale.verseMeaning)}</text>
    <text x="${locale.rtl ? x + pad : x + w - pad}" y="${y + h * 0.77}" text-anchor="${locale.rtl ? "start" : "end"}" fill="#a9ada8" font-family="Arial, sans-serif" font-size="${w * 0.025}">94:5</text>
    <rect x="${x + pad}" y="${y + h * 0.84}" width="${cardW}" height="${h * 0.12}" rx="${w * 0.04}" fill="#121a15" stroke="#2d493a"/>
    <text x="${tx}" y="${y + h * 0.885}" text-anchor="${anchor}" fill="#fff8e8" font-family="Arial, sans-serif" font-size="${w * 0.03}" font-weight="700">${escapeXml(locale.prayerTitle)}</text>
    <text x="${tx}" y="${y + h * 0.925}" text-anchor="${anchor}" fill="#8f9892" font-family="Arial, sans-serif" font-size="${w * 0.024}">${escapeXml(locale.prayerPlace)}</text>
    <text x="${locale.rtl ? x + pad : x + w - pad}" y="${y + h * 0.905}" text-anchor="${locale.rtl ? "start" : "end"}" fill="#e7b72b" font-family="Arial, sans-serif" font-size="${w * 0.036}" font-weight="700">20:18</text>
  `;
}

function sceneAnswer(locale, box) {
  const { x, y, width: w, height: h } = box;
  const pad = w * 0.06;
  const contentW = w - pad * 2;
  const anchor = locale.rtl ? "end" : "start";
  const tx = locale.rtl ? x + w - pad : x + pad;
  const questionX = locale.rtl ? x + pad : x + w * 0.31;
  return `
    <text x="${tx}" y="${y + h * 0.07}" text-anchor="${anchor}" fill="#e7b72b" font-family="Arial, sans-serif" font-size="${w * 0.044}" font-weight="800">KuranAI</text>
    <g transform="translate(${x + pad} 0)">${sourceTabs(locale, contentW, y + h * 0.11, true)}</g>
    <rect x="${questionX}" y="${y + h * 0.22}" width="${w * 0.63}" height="${h * 0.085}" rx="${w * 0.045}" fill="#e7b72b"/>
    <text x="${questionX + w * 0.315}" y="${y + h * 0.272}" text-anchor="middle" fill="#0b100d" font-family="Arial, sans-serif" font-size="${w * 0.027}" font-weight="700">${escapeXml(locale.question)}</text>
    <rect x="${x + pad}" y="${y + h * 0.34}" width="${contentW}" height="${h * 0.43}" rx="${w * 0.045}" fill="#171510" stroke="#3c3525"/>
    <circle cx="${locale.rtl ? x + w - pad - w * 0.04 : x + pad + w * 0.04}" cy="${y + h * 0.39}" r="${w * 0.035}" fill="#3a3216"/>
    <text x="${locale.rtl ? x + w - pad - w * 0.04 : x + pad + w * 0.04}" y="${y + h * 0.405}" text-anchor="middle" fill="#efbe2b" font-family="Geeza Pro, Arial, sans-serif" font-size="${w * 0.044}" font-weight="700">ق</text>
    <text x="${tx}" y="${y + h * 0.46}" text-anchor="${anchor}" fill="#fff8e8" font-family="Arial, sans-serif" font-size="${w * 0.04}" font-weight="800">${escapeXml(locale.answerTitle)}</text>
    ${textBlock({
      lines: locale.answerLines,
      x: tx,
      y: y + h * 0.525,
      fontSize: w * 0.029,
      lineHeight: h * 0.045,
      fill: "#d9d8d1",
      weight: 400,
      anchor,
    })}
    <text x="${tx}" y="${y + h * 0.68}" text-anchor="${anchor}" fill="#9ca29d" font-family="Arial, sans-serif" font-size="${w * 0.023}">${escapeXml(locale.sources)}</text>
    <rect x="${x + pad}" y="${y + h * 0.705}" width="${w * 0.25}" height="${h * 0.045}" rx="${h * 0.0225}" fill="#332b12"/>
    <text x="${x + pad + w * 0.125}" y="${y + h * 0.735}" text-anchor="middle" fill="#efbe2b" font-family="Arial, sans-serif" font-size="${w * 0.022}" font-weight="700">2:153</text>
    <rect x="${x + pad + w * 0.27}" y="${y + h * 0.705}" width="${w * 0.25}" height="${h * 0.045}" rx="${h * 0.0225}" fill="#332b12"/>
    <text x="${x + pad + w * 0.395}" y="${y + h * 0.735}" text-anchor="middle" fill="#efbe2b" font-family="Arial, sans-serif" font-size="${w * 0.022}" font-weight="700">3:200</text>
    <rect x="${x + pad}" y="${y + h * 0.84}" width="${contentW}" height="${h * 0.075}" rx="${h * 0.0375}" fill="#181713" stroke="#3b382f"/>
    <text x="${tx}" y="${y + h * 0.888}" text-anchor="${anchor}" fill="#777b77" font-family="Arial, sans-serif" font-size="${w * 0.026}">${escapeXml(locale.askPlaceholder)}</text>
    ${micIcon(locale.rtl ? x + pad + w * 0.055 : x + w - pad - w * 0.055, y + h * 0.877, w * 0.06, "answer")}
  `;
}

function sceneVoice(locale, box, id) {
  const { x, y, width: w, height: h } = box;
  const centerX = x + w / 2;
  const waveY = y + h * 0.62;
  const bars = Array.from({ length: 17 }, (_, index) => {
    const distance = Math.abs(index - 8);
    const barHeight = h * (0.025 + (8 - distance) * 0.007);
    return `<rect x="${centerX - w * 0.34 + index * w * 0.042}" y="${waveY - barHeight / 2}" width="${w * 0.014}" height="${barHeight}" rx="${w * 0.007}" fill="${index % 2 ? "#e7b72b" : "#fff0a6"}" opacity="${0.55 + (8 - distance) * 0.05}"/>`;
  }).join("");
  return `
    <text x="${centerX}" y="${y + h * 0.11}" text-anchor="middle" fill="#e7b72b" font-family="Arial, sans-serif" font-size="${w * 0.045}" font-weight="800">KuranAI</text>
    <text x="${centerX}" y="${y + h * 0.18}" text-anchor="middle" fill="#fff8e8" font-family="Arial, sans-serif" font-size="${w * 0.036}" font-weight="700">${escapeXml(locale.voiceHint)}</text>
    <circle cx="${centerX}" cy="${y + h * 0.42}" r="${w * 0.22}" fill="#d8a613" opacity=".09"/>
    <circle cx="${centerX}" cy="${y + h * 0.42}" r="${w * 0.16}" fill="#d8a613" opacity=".13"/>
    ${micIcon(centerX, y + h * 0.42, w * 0.19, id)}
    ${bars}
    <text x="${centerX}" y="${y + h * 0.71}" text-anchor="middle" fill="#fff8e8" font-family="Arial, sans-serif" font-size="${w * 0.043}" font-weight="800">${escapeXml(locale.listening)}</text>
    <rect x="${x + w * 0.12}" y="${y + h * 0.77}" width="${w * 0.76}" height="${h * 0.085}" rx="${h * 0.0425}" fill="#171510" stroke="#4a3f24"/>
    <circle cx="${x + w * 0.18}" cy="${y + h * 0.812}" r="${w * 0.025}" fill="#e7b72b"/>
    <path d="M ${x + w * 0.168} ${y + h * 0.812} l ${w * 0.008} ${h * 0.009} l ${w * 0.018} ${-h * 0.02}" fill="none" stroke="#07130f" stroke-width="${w * 0.008}" stroke-linecap="round"/>
    <text x="${centerX + (locale.rtl ? w * 0.05 : 0)}" y="${y + h * 0.822}" text-anchor="middle" fill="#d9d8d1" font-family="Arial, sans-serif" font-size="${w * 0.027}" font-weight="600">${escapeXml(locale.autoStop)}</text>
  `;
}

function sceneLanguages(locale, box) {
  const { x, y, width: w, height: h } = box;
  const centerX = x + w / 2;
  const pad = w * 0.07;
  const rows = locale.historyRows
    .map((row, index) => {
      const rowY = y + h * (0.47 + index * 0.13);
      const anchor = locale.rtl ? "end" : "start";
      const tx = locale.rtl ? x + w - pad - w * 0.04 : x + pad + w * 0.04;
      return `<rect x="${x + pad}" y="${rowY}" width="${w - pad * 2}" height="${h * 0.095}" rx="${w * 0.035}" fill="#171510" stroke="#353126"/>
        <circle cx="${locale.rtl ? x + w - pad - w * 0.04 : x + pad + w * 0.04}" cy="${rowY + h * 0.047}" r="${w * 0.026}" fill="#3a3216"/>
        <text x="${tx + (locale.rtl ? -w * 0.055 : w * 0.055)}" y="${rowY + h * 0.055}" text-anchor="${anchor}" fill="#ece8dc" font-family="Arial, sans-serif" font-size="${w * 0.025}" font-weight="600">${escapeXml(row)}</text>`;
    })
    .join("");
  return `
    <text x="${centerX}" y="${y + h * 0.1}" text-anchor="middle" fill="#e7b72b" font-family="Arial, sans-serif" font-size="${w * 0.043}" font-weight="800">KuranAI</text>
    <text x="${centerX}" y="${y + h * 0.17}" text-anchor="middle" fill="#fff8e8" font-family="Arial, sans-serif" font-size="${w * 0.036}" font-weight="700">${escapeXml(locale.languageTitle)}</text>
    <g transform="translate(${x + w * 0.1} ${y + h * 0.23})">
      <rect width="${w * 0.8}" height="${h * 0.14}" rx="${h * 0.07}" fill="#171510" stroke="#4a3f24"/>
      <rect x="${w * 0.025}" y="${h * 0.024}" width="${w * 0.235}" height="${h * 0.092}" rx="${h * 0.046}" fill="#e7b72b"/>
      <text x="${w * 0.142}" y="${h * 0.083}" text-anchor="middle" fill="#07130f" font-family="Arial, sans-serif" font-size="${w * 0.027}" font-weight="800">Türkçe</text>
      <text x="${w * 0.4}" y="${h * 0.083}" text-anchor="middle" fill="#ece8dc" font-family="Arial, sans-serif" font-size="${w * 0.027}" font-weight="700">English</text>
      <text x="${w * 0.665}" y="${h * 0.083}" text-anchor="middle" fill="#ece8dc" font-family="Geeza Pro, Arial, sans-serif" font-size="${w * 0.03}" font-weight="700">العربية</text>
    </g>
    <text x="${locale.rtl ? x + w - pad : x + pad}" y="${y + h * 0.43}" text-anchor="${locale.rtl ? "end" : "start"}" fill="#9ca29d" font-family="Arial, sans-serif" font-size="${w * 0.025}" font-weight="700">${escapeXml(locale.historyTitle)}</text>
    ${rows}
  `;
}

function screenshotSvg({ width, height, localeKey, sceneIndex, platform }) {
  const locale = locales[localeKey];
  const slide = locale.slides[sceneIndex];
  const id = `${platform}-${localeKey}-${sceneIndex}`.replaceAll(/[^a-zA-Z0-9-]/g, "");
  const margin = width * 0.075;
  const anchor = locale.rtl ? "end" : "start";
  const titleX = locale.rtl ? width - margin : margin;
  const titleY = height * 0.085;
  const titleSize = width * (platform === "app-store" ? 0.085 : 0.078);
  const titleLine = titleSize * 1.04;
  const phone = {
    x: width * 0.1,
    y: height * (platform === "app-store" ? 0.27 : 0.285),
    width: width * 0.8,
    height: height * (platform === "app-store" ? 0.79 : 0.78),
  };
  const scenes = [
    () => sceneHome(locale, phone, id),
    () => sceneAnswer(locale, phone),
    () => sceneVoice(locale, phone, id),
    () => sceneLanguages(locale, phone),
  ];
  return `<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
    ${baseDefs(id)}
    <rect width="${width}" height="${height}" fill="url(#bg-${id})"/>
    <rect width="${width}" height="${height}" fill="url(#glow-${id})"/>
    <circle cx="${locale.rtl ? width * 0.13 : width * 0.87}" cy="${height * 0.08}" r="${width * 0.22}" fill="#c69a1f" opacity=".05"/>
    <text x="${titleX}" y="${height * 0.045}" text-anchor="${anchor}" fill="#e7b72b" font-family="Arial, sans-serif" font-size="${width * 0.022}" font-weight="800" letter-spacing="${locale.rtl ? 0 : width * 0.004}">${escapeXml(locale.eyebrow)}</text>
    ${textBlock({
      lines: slide.title,
      x: titleX,
      y: titleY,
      fontSize: titleSize,
      lineHeight: titleLine,
      fill: "#fff8e8",
      weight: 800,
      anchor,
      family: locale.rtl ? "Geeza Pro, Arial, sans-serif" : "Arial, sans-serif",
    })}
    ${textBlock({
      lines: [slide.subtitle],
      x: titleX,
      y: titleY + titleLine * 2.25,
      fontSize: width * 0.029,
      lineHeight: width * 0.04,
      fill: "#aeb5af",
      weight: 400,
      anchor,
      family: locale.rtl ? "Geeza Pro, Arial, sans-serif" : "Arial, sans-serif",
    })}
    ${phoneChrome({ ...phone, id, content: scenes[sceneIndex]() })}
  </svg>`;
}

function featureGraphicSvg(localeKey) {
  const locale = locales[localeKey];
  const width = 1024;
  const height = 500;
  const anchor = locale.rtl ? "end" : "start";
  const tx = locale.rtl ? 955 : 69;
  const logoX = locale.rtl ? 105 : 919;
  const id = `feature-${localeKey}`.replaceAll(/[^a-zA-Z0-9-]/g, "");
  return `<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
    ${baseDefs(id)}
    <rect width="${width}" height="${height}" fill="url(#bg-${id})"/>
    <rect width="${width}" height="${height}" fill="url(#glow-${id})"/>
    <circle cx="${logoX}" cy="250" r="150" fill="#d2a322" opacity=".07"/>
    <rect x="${logoX - 92}" y="158" width="184" height="184" rx="45" fill="#f2bd24"/>
    <text x="${logoX}" y="284" text-anchor="middle" fill="#07130f" font-family="Geeza Pro, Arial, sans-serif" font-size="112" font-weight="700">ق</text>
    <text x="${tx}" y="78" text-anchor="${anchor}" fill="#e7b72b" font-family="Arial, sans-serif" font-size="27" font-weight="800" letter-spacing="${locale.rtl ? 0 : 4}">KURANAI</text>
    ${textBlock({
      lines: locale.featureTitle,
      x: tx,
      y: 154,
      fontSize: locale.rtl ? 48 : 50,
      lineHeight: 58,
      fill: "#fff8e8",
      weight: 800,
      anchor,
      family: locale.rtl ? "Geeza Pro, Arial, sans-serif" : "Arial, sans-serif",
    })}
    <text x="${tx}" y="302" text-anchor="${anchor}" fill="#b8c0ba" font-family="${locale.rtl ? "Geeza Pro, Arial" : "Arial"}, sans-serif" font-size="24">${escapeXml(locale.featureSubtitle)}</text>
    <rect x="${locale.rtl ? 631 : 69}" y="354" width="324" height="60" rx="30" fill="#171d18" stroke="#655321"/>
    <text x="${locale.rtl ? 793 : 231}" y="393" text-anchor="middle" fill="#f2bd24" font-family="Arial, sans-serif" font-size="21" font-weight="700">Türkçe  •  English  •  العربية</text>
  </svg>`;
}

async function renderSvg(svg, svgPath, pngPath, width, height) {
  const normalizedSvg = `${svg.replace(/[ \t]+$/gm, "")}\n`;
  await writeFile(svgPath, normalizedSvg, "utf8");
  await run("/opt/homebrew/bin/rsvg-convert", [
    "--width",
    String(width),
    "--height",
    String(height),
    "--format",
    "png",
    "--output",
    pngPath,
    svgPath,
  ]);
}

for (const [platform, spec] of Object.entries(specs)) {
  for (const localeKey of Object.keys(locales)) {
    const output = path.join(root, platform, localeKey);
    await mkdir(output, { recursive: true });
    for (let sceneIndex = 0; sceneIndex < 4; sceneIndex += 1) {
      const baseName = `${sceneIndex + 1}-${["quran", "sources", "voice", "languages"][sceneIndex]}`;
      await renderSvg(
        screenshotSvg({ ...spec, localeKey, sceneIndex, platform }),
        path.join(output, `${baseName}.svg`),
        path.join(output, `${baseName}.png`),
        spec.width,
        spec.height,
      );
    }
  }
}

for (const localeKey of Object.keys(locales)) {
  const output = path.join(root, "google-play", localeKey);
  await renderSvg(
    featureGraphicSvg(localeKey),
    path.join(output, "feature-graphic.svg"),
    path.join(output, "feature-graphic.png"),
    1024,
    500,
  );
}

console.log(`Generated localized ASO assets in ${root}`);
