package com.talhacan.kuranai.model

import org.junit.Assert.assertEquals
import org.junit.Test

class ModelsTest {
    @Test
    fun languageCodesRoundTripAndUnknownDefaultsToTurkish() {
        AppLanguage.entries.forEach { language ->
            assertEquals(language, AppLanguage.from(language.code))
        }
        assertEquals(AppLanguage.Turkish, AppLanguage.from("de"))
        assertEquals(AppLanguage.Turkish, AppLanguage.from(null))
    }

    @Test
    fun packageNamesAreLocalizedForEverySupportedLanguage() {
        val subscriptionPackage = SubscriptionPackage(
            id = 1,
            name = "Altın",
            nameEn = "Gold",
            slug = "gold",
            dailyQuestionLimit = 50,
            dailyTtsLimit = 10,
            monthlyPriceTl = 79.99,
            yearlyPriceTl = 799.99,
        )

        assertEquals("Altın", subscriptionPackage.localizedName(AppLanguage.Turkish))
        assertEquals("Gold", subscriptionPackage.localizedName(AppLanguage.English))
        assertEquals("ذهبية", subscriptionPackage.localizedName(AppLanguage.Arabic))
    }
}
