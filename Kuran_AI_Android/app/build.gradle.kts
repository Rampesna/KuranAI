import java.util.Properties
import org.gradle.api.GradleException

plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.compose)
}

android {
    namespace = "com.talhacan.kuranai"
    compileSdk {
        version = release(36)
    }

    defaultConfig {
        applicationId = "com.talhacan.kuranai"
        minSdk = 28
        targetSdk = 36
        versionCode = 4
        versionName = "1.0.3"

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
        buildConfigField("String", "API_BASE_URL", "\"https://kuranai-api.talhacan.com\"")
        buildConfigField(
            "String",
            "GOOGLE_WEB_CLIENT_ID",
            "\"157832066999-3tgb2daavmrd1d1vd10ptrato1iu382t.apps.googleusercontent.com\"",
        )
        buildConfigField(
            "String",
            "APPLE_SIGN_IN_URL",
            "\"https://kuranai.talhacan.com/auth/apple/android\"",
        )
    }

    val releaseSigningRequested = gradle.startParameter.taskNames.any {
        it.contains("release", ignoreCase = true)
    }
    var releaseSigningAvailable = false

    signingConfigs {
        create("release") {
            val propertiesFile = rootProject.file("keystore.properties")
            if (propertiesFile.exists()) {
                val properties = Properties().apply {
                    propertiesFile.inputStream().use(::load)
                }
                storeFile = file(properties.getProperty("storeFile"))
                storePassword = properties.getProperty("storePassword")
                keyAlias = properties.getProperty("keyAlias")
                keyPassword = properties.getProperty("keyPassword")
                releaseSigningAvailable = true
            } else {
                val environment = System.getenv()
                val storePath = environment["KURANAI_UPLOAD_KEYSTORE_PATH"]
                val storeSecret = environment["KURANAI_UPLOAD_STORE_PASSWORD"]
                val alias = environment["KURANAI_UPLOAD_KEY_ALIAS"]
                val keySecret = environment["KURANAI_UPLOAD_KEY_PASSWORD"]
                if (
                    !storePath.isNullOrBlank() &&
                    !storeSecret.isNullOrBlank() &&
                    !alias.isNullOrBlank() &&
                    !keySecret.isNullOrBlank()
                ) {
                    storeFile = file(storePath)
                    storePassword = storeSecret
                    keyAlias = alias
                    keyPassword = keySecret
                    releaseSigningAvailable = true
                } else {
                    storeFile = file("${System.getProperty("user.home")}/.android/debug.keystore")
                    storePassword = "android"
                    keyAlias = "androiddebugkey"
                    keyPassword = "android"
                }
            }
        }
    }

    if (releaseSigningRequested && !releaseSigningAvailable) {
        throw GradleException(
            "Release signing credentials are required. Configure keystore.properties " +
                "or the KURANAI_UPLOAD_* environment variables.",
        )
    }

    buildTypes {
        debug {
            applicationIdSuffix = ".debug"
            versionNameSuffix = "-debug"
        }
        release {
            isDebuggable = false
            isMinifyEnabled = false
            signingConfig = signingConfigs.getByName("release")
        }
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }
    buildFeatures {
        compose = true
        buildConfig = true
    }
    bundle {
        language {
            enableSplit = false
        }
    }
    packaging {
        resources.excludes += "/META-INF/{AL2.0,LGPL2.1}"
    }
}

dependencies {
    implementation(platform(libs.androidx.compose.bom))
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.activity.compose)
    implementation(libs.androidx.lifecycle.runtime.ktx)
    implementation(libs.androidx.lifecycle.viewmodel.compose)
    implementation(libs.androidx.compose.ui)
    implementation(libs.androidx.compose.ui.graphics)
    implementation(libs.androidx.compose.ui.tooling.preview)
    implementation(libs.androidx.compose.material3)
    implementation(libs.androidx.compose.material.icons.extended)
    implementation(libs.androidx.navigation.compose)
    implementation(libs.kotlinx.coroutines.android)
    implementation(libs.okhttp)
    implementation(libs.okhttp.logging)
    implementation(libs.coil.compose)
    implementation(libs.androidx.credentials)
    implementation(libs.androidx.credentials.play.services)
    implementation(libs.googleid)
    implementation(libs.billing.ktx)
    implementation(libs.androidx.browser)
    implementation(libs.androidx.security.crypto)
    implementation(libs.gson)
    implementation(libs.play.review)
    implementation(libs.play.review.ktx)

    testImplementation(libs.junit)
    androidTestImplementation(platform(libs.androidx.compose.bom))
    androidTestImplementation(libs.androidx.compose.ui.test.junit4)
    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso)
    debugImplementation(libs.androidx.compose.ui.tooling)
    debugImplementation(libs.androidx.compose.ui.test.manifest)
}
