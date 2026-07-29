package com.talhacan.kuranai.ui

import com.talhacan.kuranai.model.AppLanguage

fun AppLanguage.t(turkish: String, english: String, arabic: String): String = when (this) {
    AppLanguage.Turkish -> turkish
    AppLanguage.English -> english
    AppLanguage.Arabic -> arabic
}

fun AppLanguage.localeTag(): String = when (this) {
    AppLanguage.Turkish -> "tr-TR"
    AppLanguage.English -> "en-US"
    AppLanguage.Arabic -> "ar-SA"
}
