const fs = require('fs');
const path = require('path');

const sqlPath = path.join(__dirname, 'myself3dblog_2026-10-05_01-13-05_mysql_data_vw18p.sql');
const content = fs.readFileSync(sqlPath, 'utf8');

// Find INSERT INTO `articles`
const articleInsertsMatch = content.match(/INSERT INTO `articles` VALUES\s*([\s\S]+?);(?=\s*(?:INSERT|UNLOCK|DROP|CREATE|ALTER|$))/);
if (!articleInsertsMatch) {
  console.log('No articles insert found');
  process.exit(1);
}

const rawValues = articleInsertsMatch[1];
console.log('Raw values length:', rawValues.length);

// Custom SQL row parser that respects single quotes and escape sequences
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
          // Check if double quote escape in SQL
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
console.log(`Parsed ${rows.length} articles!`);
rows.slice(0, 5).forEach((r, idx) => {
  console.log(`[#${idx + 1}] ID:${r[0]} | Title:${r[1]} | TitleEn:${r[2]} | Cat:${r[3]} | Price:${r[11]} | Cover:${r[8]}`);
});
