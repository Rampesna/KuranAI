package com.talhacan.kuranai

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.viewModels
import com.talhacan.kuranai.ui.KuranAIApp

class MainActivity : ComponentActivity() {
    private val viewModel by viewModels<KuranAIViewModel>()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        handleIntent(intent)
        setContent {
            KuranAIApp(viewModel = viewModel, activity = this)
        }
    }

    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        setIntent(intent)
        handleIntent(intent)
    }

    private fun handleIntent(intent: Intent?) {
        val uri = intent?.data ?: return
        if (uri.scheme == "kuranai" && uri.host == "auth" && uri.path == "/apple") {
            uri.getQueryParameter("code")
                ?.takeIf(String::isNotBlank)
                ?.let(viewModel::completeAppleSignIn)
        }
    }
}
