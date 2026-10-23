// WallArt Encrypted Asset Decryptor
// Matches Go backend WallArtSecretKey = "WallArtTendiesKey2026"
const WALLART_KEY = 'WallArtTendiesKey2026'

// In-memory cache: url -> decrypted Blob URL (blob:http://...)
const blobCache = new Map()

/**
 * Fetch and decrypt an encrypted .dat asset URL in memory, returning a blob: URL.
 * If the URL is unencrypted (.png, .jpg) or already decrypted, returns it immediately.
 */
export async function resolveAndDecryptAssetUrl(url) {
  if (!url) return ''

  // If not a .dat encrypted asset, return as is
  if (!url.includes('.dat')) {
    return url
  }

  if (blobCache.has(url)) {
    return blobCache.get(url)
  }

  try {
    const res = await fetch(url)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const arrayBuffer = await res.arrayBuffer()
    const bytes = new Uint8Array(arrayBuffer)

    // Check WART magic header: 0x57, 0x41, 0x52, 0x54 ('W', 'A', 'R', 'T')
    if (
      bytes.length >= 6 &&
      bytes[0] === 0x57 &&
      bytes[1] === 0x41 &&
      bytes[2] === 0x52 &&
      bytes[3] === 0x54
    ) {
      const salt = bytes[5]
      const rawLen = bytes.length - 6
      const decrypted = new Uint8Array(rawLen)
      const keyLen = WALLART_KEY.length

      for (let i = 0; i < rawLen; i++) {
        const k = WALLART_KEY.charCodeAt((i + salt) % keyLen)
        const shift = (i * 7 + salt) & 0xff
        decrypted[i] = bytes[6 + i] ^ k ^ shift
      }

      // Detect MIME type from decrypted bytes
      let mimeType = 'image/png'
      if (decrypted[0] === 0xff && decrypted[1] === 0xd8) {
        mimeType = 'image/jpeg'
      } else if (
        decrypted[0] === 0x52 &&
        decrypted[1] === 0x49 &&
        decrypted[2] === 0x46 &&
        decrypted[3] === 0x46
      ) {
        mimeType = 'image/webp'
      }

      const blob = new Blob([decrypted], { type: mimeType })
      const blobUrl = URL.createObjectURL(blob)
      blobCache.set(url, blobUrl)
      return blobUrl
    } else {
      // Unencrypted, fallback to original
      return url
    }
  } catch (err) {
    console.warn('[WallArtCrypto] Decrypt failed, fallback to raw url:', err)
    return url
  }
}
