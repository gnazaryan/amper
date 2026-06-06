import Moment from 'moment';

export const formatDate = (input) => {
  if (input != null) {
    const date = new Date(input);
    return isToday(date) ? 'Today ' + Moment(date).format('h:mm:ss a') : isYesterday(date) ? 'Yesterday ' + Moment(date).format('h:mm:ss a') : Moment(date).format('MMMM Do YYYY, h:mm:ss a');
  }
  return null;
};

const isToday = (someDate) => {
  const today = new Date()
  return someDate.getDate() == today.getDate() &&
    someDate.getMonth() == today.getMonth() &&
    someDate.getFullYear() == today.getFullYear()
}

const isYesterday = (someDate) => {
  const today = new Date()
  const yesterday = today.getDate() - 1;
  return someDate.getDate() == yesterday &&
    someDate.getMonth() == today.getMonth() &&
    someDate.getFullYear() == today.getFullYear()
}

export const copyArray = (input) => {
  if (input != null) {
    const result = [];
    for (let i = 0; i < input.length; i++) {
      result[i] = input[i];
    }
    return result;
  }
  return null;
};

export const truncate = ( str, n, useWordBoundary ) => {
  if (str.length <= n) { return str; }
  const subString = str.slice(0, n-1); // the original check
  return (useWordBoundary 
    ? subString.slice(0, subString.lastIndexOf(" ")) 
    : subString) + "...";
};

export const formatVersion = (version) => {
  return `${version.major}.${version.minor}.${version.patch}`
};

export const uniqueId = () => {
  var uniqid = Date.now() + '';
  var bytes = [];
  for (var i = 0; i < uniqid.length; ++i) {
    bytes = bytes.concat([uniqid.charCodeAt(i)]);
  }
  return base64EncArr(bytes);
};

export const setStoreValue = (key, value) => {
  if (typeof(Storage) !== "undefined") {
    localStorage.setItem(key, value);
  }
};

export const gettStoreValue = (key) => {
  if (typeof(Storage) !== "undefined") {
    return localStorage.getItem(key);
  }
};

export const makeApieName = (value) => {
    let apiName = value.toLowerCase();
    let result = '';
    let spaceFound = false;
    for(let i = 0; i < apiName.length; i++) {
        const character = apiName.charAt(i);
        if (character.match(/[a-z]/i)) {
            result = result + (spaceFound ? character.toUpperCase(): character);
            spaceFound = false;
        } else if(character === ' ') {
            spaceFound = true;
        }
    }
    if (result.length > 0) {
        result = result + '_amp';
    }
    return result;
}

export const registerResize = (callback, option) => {
    const resize = (func, option) => {
        let timer;
        return function(event) {
            if(timer) {
                clearTimeout(timer);
            }
            timer = setTimeout(func, 1000, event, option);
        };
    };
    window.addEventListener('resize', resize(function(event, option) {
        callback(window.innerHeight, window.innerWidth, option);
    }, option));
}

export const clone = (input) => {
    return JSON.parse(JSON.stringify(input));
}

export const requireExternal = (src) => {
  const script = document.createElement("script");
  script.src = src;
  script.async = true;
  document.body.appendChild(script);
}

export const replace = (items, item) => {
    let indexOf = -1;
    for (let i = 0; i < items.length; i++) {
        if (items[i].id === item.id) {
            indexOf = i;
            break;
        }
    }
    
    if (indexOf > -1) {
        items[indexOf] = item;
    } else {
        items.push(item);
    }
}

export const debounce = (callback, time) => {
    let timer;
    let items = [];
    const getItems = () => {
        const temp = items;
        items = [];
        return temp;
    };
    return (item) => {
        if(timer) {
            clearTimeout(timer);
        }
        replace(items, item);
        timer = setTimeout(callback, time || 1000, getItems);
    };
};

export const debounceLatest = (callback, time) => {
    let timer;
    let item = null;
    const getItem = () => {
        return item;
    };
    return (newItem, option) => {
        if(timer) {
            clearTimeout(timer);
        }
        item = newItem;
        timer = setTimeout(callback, time || 1000, getItem, option);
    };
};

  /* Base64 string to array encoding */
  export const uint6ToB64 = (nUint6) => {
    return nUint6 < 26
      ? nUint6 + 65
      : nUint6 < 52
      ? nUint6 + 71
      : nUint6 < 62
      ? nUint6 - 4
      : nUint6 === 62
      ? 43
      : nUint6 === 63
      ? 47
      : 65;
  }

  // Array of bytes to Base64 string decoding
  export const b64ToUint6 = (nChr) => {
    return nChr > 64 && nChr < 91
      ? nChr - 65
      : nChr > 96 && nChr < 123
      ? nChr - 71
      : nChr > 47 && nChr < 58
      ? nChr + 4
      : nChr === 43
      ? 62
      : nChr === 47
      ? 63
      : 0;
  }
export const base64DecToArr = (sBase64, nBlocksSize) => {
    const sB64Enc = sBase64.replace(/[^A-Za-z0-9+/]/g, ""); // Only necessary if the base64 includes whitespace such as line breaks.
    const nInLen = sB64Enc.length;
    const nOutLen = nBlocksSize
      ? Math.ceil(((nInLen * 3 + 1) >> 2) / nBlocksSize) * nBlocksSize
      : (nInLen * 3 + 1) >> 2;
    const taBytes = new Uint8Array(nOutLen);
  
    let nMod3;
    let nMod4;
    let nUint24 = 0;
    let nOutIdx = 0;
    for (let nInIdx = 0; nInIdx < nInLen; nInIdx++) {
      nMod4 = nInIdx & 3;
      nUint24 |= b64ToUint6(sB64Enc.charCodeAt(nInIdx)) << (6 * (3 - nMod4));
      if (nMod4 === 3 || nInLen - nInIdx === 1) {
        nMod3 = 0;
        while (nMod3 < 3 && nOutIdx < nOutLen) {
          taBytes[nOutIdx] = (nUint24 >>> ((16 >>> nMod3) & 24)) & 255;
          nMod3++;
          nOutIdx++;
        }
        nUint24 = 0;
      }
    }
  
    return taBytes;
  };
  
  export const base64EncArr = (aBytes) => {
    let nMod3 = 2;
    let sB64Enc = "";
  
    const nLen = aBytes.length;
    let nUint24 = 0;
    for (let nIdx = 0; nIdx < nLen; nIdx++) {
      nMod3 = nIdx % 3;
      // To break your base64 into several 80-character lines, add:
      //   if (nIdx > 0 && ((nIdx * 4) / 3) % 76 === 0) {
      //      sB64Enc += "\r\n";
      //    }
  
      nUint24 |= aBytes[nIdx] << ((16 >>> nMod3) & 24);
      if (nMod3 === 2 || aBytes.length - nIdx === 1) {
        sB64Enc += String.fromCodePoint(
          uint6ToB64((nUint24 >>> 18) & 63),
          uint6ToB64((nUint24 >>> 12) & 63),
          uint6ToB64((nUint24 >>> 6) & 63),
          uint6ToB64(nUint24 & 63)
        );
        nUint24 = 0;
      }
    }
    return (
      sB64Enc.substring(0, sB64Enc.length - 2 + nMod3) +
      (nMod3 === 2 ? "" : nMod3 === 1 ? "=" : "==")
    );
  }
  
  /* UTF-8 array to JS string and vice versa */
  
  export const UTF8ArrToStr = (aBytes) => {
    let sView = "";
    let nPart;
    const nLen = aBytes.length;
    for (let nIdx = 0; nIdx < nLen; nIdx++) {
      nPart = aBytes[nIdx];
      sView += String.fromCodePoint(
        nPart > 251 && nPart < 254 && nIdx + 5 < nLen /* six bytes */
          ? /* (nPart - 252 << 30) may be not so safe in ECMAScript! So…: */
            (nPart - 252) * 1073741824 +
              ((aBytes[++nIdx] - 128) << 24) +
              ((aBytes[++nIdx] - 128) << 18) +
              ((aBytes[++nIdx] - 128) << 12) +
              ((aBytes[++nIdx] - 128) << 6) +
              aBytes[++nIdx] -
              128
          : nPart > 247 && nPart < 252 && nIdx + 4 < nLen /* five bytes */
          ? ((nPart - 248) << 24) +
            ((aBytes[++nIdx] - 128) << 18) +
            ((aBytes[++nIdx] - 128) << 12) +
            ((aBytes[++nIdx] - 128) << 6) +
            aBytes[++nIdx] -
            128
          : nPart > 239 && nPart < 248 && nIdx + 3 < nLen /* four bytes */
          ? ((nPart - 240) << 18) +
            ((aBytes[++nIdx] - 128) << 12) +
            ((aBytes[++nIdx] - 128) << 6) +
            aBytes[++nIdx] -
            128
          : nPart > 223 && nPart < 240 && nIdx + 2 < nLen /* three bytes */
          ? ((nPart - 224) << 12) +
            ((aBytes[++nIdx] - 128) << 6) +
            aBytes[++nIdx] -
            128
          : nPart > 191 && nPart < 224 && nIdx + 1 < nLen /* two bytes */
          ? ((nPart - 192) << 6) + aBytes[++nIdx] - 128
          : /* nPart < 127 ? */ /* one byte */
            nPart
      );
    }
    return sView;
  };
  
  export const strToUTF8Arr = (sDOMStr) => {
    let aBytes;
    let nChr;
    const nStrLen = sDOMStr.length;
    let nArrLen = 0;
  
    /* mapping… */
    for (let nMapIdx = 0; nMapIdx < nStrLen; nMapIdx++) {
      nChr = sDOMStr.codePointAt(nMapIdx);
  
      if (nChr >= 0x10000) {
        nMapIdx++;
      }
  
      nArrLen +=
        nChr < 0x80
          ? 1
          : nChr < 0x800
          ? 2
          : nChr < 0x10000
          ? 3
          : nChr < 0x200000
          ? 4
          : nChr < 0x4000000
          ? 5
          : 6;
    }
  
    aBytes = new Uint8Array(nArrLen);
  
    /* transcription… */
    let nIdx = 0;
    let nChrIdx = 0;
    while (nIdx < nArrLen) {
      nChr = sDOMStr.codePointAt(nChrIdx);
      if (nChr < 128) {
        /* one byte */
        aBytes[nIdx++] = nChr;
      } else if (nChr < 0x800) {
        /* two bytes */
        aBytes[nIdx++] = 192 + (nChr >>> 6);
        aBytes[nIdx++] = 128 + (nChr & 63);
      } else if (nChr < 0x10000) {
        /* three bytes */
        aBytes[nIdx++] = 224 + (nChr >>> 12);
        aBytes[nIdx++] = 128 + ((nChr >>> 6) & 63);
        aBytes[nIdx++] = 128 + (nChr & 63);
      } else if (nChr < 0x200000) {
        /* four bytes */
        aBytes[nIdx++] = 240 + (nChr >>> 18);
        aBytes[nIdx++] = 128 + ((nChr >>> 12) & 63);
        aBytes[nIdx++] = 128 + ((nChr >>> 6) & 63);
        aBytes[nIdx++] = 128 + (nChr & 63);
        nChrIdx++;
      } else if (nChr < 0x4000000) {
        /* five bytes */
        aBytes[nIdx++] = 248 + (nChr >>> 24);
        aBytes[nIdx++] = 128 + ((nChr >>> 18) & 63);
        aBytes[nIdx++] = 128 + ((nChr >>> 12) & 63);
        aBytes[nIdx++] = 128 + ((nChr >>> 6) & 63);
        aBytes[nIdx++] = 128 + (nChr & 63);
        nChrIdx++;
      } /* if (nChr <= 0x7fffffff) */ else {
        /* six bytes */
        aBytes[nIdx++] = 252 + (nChr >>> 30);
        aBytes[nIdx++] = 128 + ((nChr >>> 24) & 63);
        aBytes[nIdx++] = 128 + ((nChr >>> 18) & 63);
        aBytes[nIdx++] = 128 + ((nChr >>> 12) & 63);
        aBytes[nIdx++] = 128 + ((nChr >>> 6) & 63);
        aBytes[nIdx++] = 128 + (nChr & 63);
        nChrIdx++;
      }
      nChrIdx++;
    }
  
    return aBytes;
  };

  export const DecimalPrecision = (function() {
    if (Math.trunc === undefined) {
        Math.trunc = function(v) {
            return v < 0 ? Math.ceil(v) : Math.floor(v);
        };
    }
    var decimalAdjust = function myself(type, num, decimalPlaces) {
        if (type === 'round' && num < 0)
            return -myself(type, -num, decimalPlaces);
        var shift = function(value, exponent) {
            value = (value + 'e').split('e');
            return +(value[0] + 'e' + (+value[1] + (exponent || 0)));
        };
        var n = shift(num, +decimalPlaces);
        return shift(Math[type](n), -decimalPlaces);
    };
    return {
        // Decimal round (half away from zero)
        round: function(num, decimalPlaces) {
            return decimalAdjust('round', num, decimalPlaces);
        },
        // Decimal ceil
        ceil: function(num, decimalPlaces) {
            return decimalAdjust('ceil', num, decimalPlaces);
        },
        // Decimal floor
        floor: function(num, decimalPlaces) {
            return decimalAdjust('floor', num, decimalPlaces);
        },
        // Decimal trunc
        trunc: function(num, decimalPlaces) {
            return decimalAdjust('trunc', num, decimalPlaces);
        },
        // Format using fixed-point notation
        toFixed: function(num, decimalPlaces) {
            return decimalAdjust('round', num, decimalPlaces).toFixed(decimalPlaces);
        }
    };
})();
