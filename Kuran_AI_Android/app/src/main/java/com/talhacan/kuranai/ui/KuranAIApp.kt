package com.talhacan.kuranai.ui

import android.Manifest
import android.app.Activity
import android.content.Intent
import android.net.Uri
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.PickVisualMediaRequest
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.navigationBarsPadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.statusBarsPadding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.automirrored.filled.ExitToApp
import androidx.compose.material.icons.automirrored.filled.Send
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.AutoStories
import androidx.compose.material.icons.filled.ChatBubble
import androidx.compose.material.icons.filled.Check
import androidx.compose.material.icons.filled.Close
import androidx.compose.material.icons.filled.DarkMode
import androidx.compose.material.icons.filled.DeleteForever
import androidx.compose.material.icons.filled.Fingerprint
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.KeyboardVoice
import androidx.compose.material.icons.filled.Language
import androidx.compose.material.icons.filled.LightMode
import androidx.compose.material.icons.filled.Lock
import androidx.compose.material.icons.filled.Menu
import androidx.compose.material.icons.filled.Mic
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.PhotoCamera
import androidx.compose.material.icons.filled.PlayArrow
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.material.icons.filled.Restore
import androidx.compose.material.icons.filled.StopCircle
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.AssistChip
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Divider
import androidx.compose.material3.DrawerValue
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilledIconButton
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.IconButtonDefaults
import androidx.compose.material3.ListItem
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalNavigationDrawer
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.OutlinedCard
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.material3.Surface
import androidx.compose.material3.Switch
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.rememberDrawerState
import androidx.compose.material3.rememberModalBottomSheetState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.CompositionLocalProvider
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalLayoutDirection
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.LayoutDirection
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import coil.compose.AsyncImage
import com.talhacan.kuranai.BuildConfig
import com.talhacan.kuranai.KuranAIViewModel
import com.talhacan.kuranai.billing.BillingStore
import com.talhacan.kuranai.model.AppLanguage
import com.talhacan.kuranai.model.AppPage
import com.talhacan.kuranai.model.AppState
import com.talhacan.kuranai.model.ChatMessage
import com.talhacan.kuranai.model.SourceType
import com.talhacan.kuranai.model.SubscriptionPackage
import com.talhacan.kuranai.review.InAppReviewCoordinator
import kotlinx.coroutines.launch
import java.text.NumberFormat
import java.util.Currency
import java.util.Locale

@Composable
fun KuranAIApp(
    viewModel: KuranAIViewModel,
    activity: Activity,
) {
    val state by viewModel.state.collectAsStateWithLifecycle()
    val snackbar = remember { SnackbarHostState() }
    val reviewCoordinator = remember(activity) { InAppReviewCoordinator(activity) }
    var wasSending by remember { mutableStateOf(false) }
    val locationPermission = rememberLauncherForActivityResult(
        ActivityResultContracts.RequestMultiplePermissions(),
    ) {
        viewModel.refreshHome()
    }

    LaunchedEffect(Unit) {
        locationPermission.launch(
            arrayOf(
                Manifest.permission.ACCESS_COARSE_LOCATION,
                Manifest.permission.ACCESS_FINE_LOCATION,
            ),
        )
    }
    LaunchedEffect(state.notice, state.error) {
        val message = state.error ?: state.notice
        if (!message.isNullOrBlank()) {
            snackbar.showSnackbar(message)
            viewModel.clearStatus()
        }
    }
    LaunchedEffect(state.sending) {
        if (
            wasSending &&
            !state.sending &&
            state.error == null &&
            state.messages.lastOrNull()?.let {
                it.role == "assistant" && !it.streaming && it.content.isNotBlank()
            } == true
        ) {
            reviewCoordinator.recordSuccessfulAnswer(activity)
        }
        wasSending = state.sending
    }

    CompositionLocalProvider(
        LocalLayoutDirection provides if (state.language == AppLanguage.Arabic) {
            LayoutDirection.Rtl
        } else {
            LayoutDirection.Ltr
        },
    ) {
        KuranAITheme(darkTheme = state.darkTheme) {
            Surface(modifier = Modifier.fillMaxSize()) {
                Box {
                    when {
                        state.restoringSession -> LoadingScreen()
                        !state.authenticated -> AuthScreen(state, viewModel, activity)
                        else -> MainContent(state, viewModel, activity)
                    }
                    SnackbarHost(
                        hostState = snackbar,
                        modifier = Modifier
                            .align(Alignment.BottomCenter)
                            .navigationBarsPadding()
                            .padding(12.dp),
                    )
                    if (state.loading) {
                        Box(
                            modifier = Modifier
                                .fillMaxSize()
                                .background(Color.Black.copy(alpha = 0.32f)),
                            contentAlignment = Alignment.Center,
                        ) {
                            CircularProgressIndicator(color = KuranGold)
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun LoadingScreen() {
    Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
        CircularProgressIndicator(color = KuranGold)
    }
}

@Composable
private fun LanguageMenu(
    language: AppLanguage,
    onLanguage: (AppLanguage) -> Unit,
) {
    var expanded by remember { mutableStateOf(false) }
    Box {
        OutlinedButton(onClick = { expanded = true }) {
            Icon(Icons.Default.Language, contentDescription = null)
            Spacer(Modifier.width(6.dp))
            Text(language.displayName)
        }
        DropdownMenu(expanded = expanded, onDismissRequest = { expanded = false }) {
            AppLanguage.entries.forEach { option ->
                DropdownMenuItem(
                    text = { Text(option.displayName) },
                    leadingIcon = {
                        if (option == language) Icon(Icons.Default.Check, contentDescription = null)
                    },
                    onClick = {
                        expanded = false
                        onLanguage(option)
                    },
                )
            }
        }
    }
}

@Composable
private fun AuthScreen(
    state: AppState,
    viewModel: KuranAIViewModel,
    activity: Activity,
) {
    var registering by remember { mutableStateOf(false) }
    var name by remember { mutableStateOf("") }
    var username by remember { mutableStateOf("") }
    var email by remember { mutableStateOf("") }
    var credential by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    val l = state.language

    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .statusBarsPadding(),
        contentPadding = PaddingValues(24.dp),
        verticalArrangement = Arrangement.spacedBy(18.dp),
    ) {
        item {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.End,
            ) {
                LanguageMenu(l, viewModel::setLanguage)
            }
        }
        item {
            Column(horizontalAlignment = Alignment.CenterHorizontally, modifier = Modifier.fillMaxWidth()) {
                QuranMark(size = 76)
                Spacer(Modifier.height(16.dp))
                Text(
                    l.t("Kuran AI", "Quran AI", "قرآن AI"),
                    style = MaterialTheme.typography.headlineLarge,
                    fontWeight = FontWeight.Bold,
                )
                Text(
                    l.t(
                        "Kur'an, tefsir ve hadis asistanınız",
                        "Your Quran, Tafsir, and Hadith assistant",
                        "مساعدك للقرآن والتفسير والحديث",
                    ),
                    color = MaterialTheme.colorScheme.onSurfaceVariant,
                    textAlign = TextAlign.Center,
                )
            }
        }
        item {
            Card(
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant),
            ) {
                Column(
                    modifier = Modifier.padding(18.dp),
                    verticalArrangement = Arrangement.spacedBy(12.dp),
                ) {
                    Button(
                        onClick = { viewModel.signInGoogle(activity) },
                        modifier = Modifier.fillMaxWidth(),
                    ) {
                        Text(l.t("Google ile devam et", "Continue with Google", "المتابعة عبر Google"))
                    }
                    OutlinedButton(
                        onClick = { viewModel.openAppleSignIn(activity) },
                        modifier = Modifier.fillMaxWidth(),
                    ) {
                        Text(l.t("Apple ile devam et", "Continue with Apple", "المتابعة عبر Apple"))
                    }
                    OutlinedButton(
                        onClick = { viewModel.signInPasskey(activity) },
                        modifier = Modifier.fillMaxWidth(),
                    ) {
                        Icon(Icons.Default.Fingerprint, contentDescription = null)
                        Spacer(Modifier.width(8.dp))
                        Text(l.t("Geçiş anahtarı ile giriş", "Sign in with a passkey", "الدخول بمفتاح مرور"))
                    }
                    HorizontalDivider()
                    Row(Modifier.fillMaxWidth()) {
                        TextButton(
                            onClick = { registering = false },
                            modifier = Modifier.weight(1f),
                        ) {
                            Text(l.t("Giriş yap", "Sign in", "تسجيل الدخول"))
                        }
                        TextButton(
                            onClick = { registering = true },
                            modifier = Modifier.weight(1f),
                        ) {
                            Text(l.t("Hesap oluştur", "Create account", "إنشاء حساب"))
                        }
                    }
                    if (registering) {
                        OutlinedTextField(
                            value = name,
                            onValueChange = { name = it },
                            label = { Text(l.t("Adınız", "Your name", "اسمك")) },
                            modifier = Modifier.fillMaxWidth(),
                            singleLine = true,
                        )
                        OutlinedTextField(
                            value = username,
                            onValueChange = { username = it },
                            label = { Text(l.t("Kullanıcı adı", "Username", "اسم المستخدم")) },
                            modifier = Modifier.fillMaxWidth(),
                            singleLine = true,
                        )
                        OutlinedTextField(
                            value = email,
                            onValueChange = { email = it },
                            label = { Text(l.t("E-posta", "Email", "البريد الإلكتروني")) },
                            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
                            modifier = Modifier.fillMaxWidth(),
                            singleLine = true,
                        )
                    } else {
                        OutlinedTextField(
                            value = credential,
                            onValueChange = { credential = it },
                            label = {
                                Text(
                                    l.t(
                                        "E-posta veya kullanıcı adı",
                                        "Email or username",
                                        "البريد الإلكتروني أو اسم المستخدم",
                                    ),
                                )
                            },
                            modifier = Modifier.fillMaxWidth(),
                            singleLine = true,
                        )
                    }
                    OutlinedTextField(
                        value = password,
                        onValueChange = { password = it },
                        label = { Text(l.t("Şifre", "Password", "كلمة المرور")) },
                        visualTransformation = PasswordVisualTransformation(),
                        keyboardOptions = KeyboardOptions(
                            keyboardType = KeyboardType.Password,
                            imeAction = ImeAction.Done,
                        ),
                        keyboardActions = KeyboardActions(
                            onDone = {
                                if (registering) {
                                    viewModel.register(name, username, email, password)
                                } else {
                                    viewModel.loginEmail(credential, password)
                                }
                            },
                        ),
                        modifier = Modifier.fillMaxWidth(),
                        singleLine = true,
                    )
                    Button(
                        onClick = {
                            if (registering) {
                                viewModel.register(name, username, email, password)
                            } else {
                                viewModel.loginEmail(credential, password)
                            }
                        },
                        enabled = password.length >= 6 &&
                            if (registering) {
                                name.isNotBlank() && username.length >= 2 && email.contains("@")
                            } else {
                                credential.isNotBlank()
                            },
                        modifier = Modifier.fillMaxWidth(),
                    ) {
                        Text(
                            if (registering) {
                                l.t("Hesabımı oluştur", "Create my account", "إنشاء حسابي")
                            } else {
                                l.t("Sohbete gir", "Enter chat", "بدء المحادثة")
                            },
                        )
                    }
                }
            }
        }
        item {
            LegalLinks(l, activity)
        }
    }
}

@Composable
private fun MainContent(
    state: AppState,
    viewModel: KuranAIViewModel,
    activity: Activity,
) {
    Scaffold(
        bottomBar = {
            if (state.page != AppPage.Plans) {
                NavigationBar {
                    NavigationBarItem(
                        selected = state.page == AppPage.Home,
                        onClick = { viewModel.setPage(AppPage.Home) },
                        icon = { Icon(Icons.Default.Home, contentDescription = null) },
                        label = {
                            Text(state.language.t("Ana sayfa", "Home", "الرئيسية"))
                        },
                    )
                    NavigationBarItem(
                        selected = state.page == AppPage.Chat,
                        onClick = { viewModel.setPage(AppPage.Chat) },
                        icon = { Icon(Icons.Default.ChatBubble, contentDescription = null) },
                        label = {
                            Text(state.language.t("Sohbet", "Chat", "المحادثة"))
                        },
                    )
                    NavigationBarItem(
                        selected = state.page == AppPage.Profile,
                        onClick = { viewModel.setPage(AppPage.Profile) },
                        icon = { Icon(Icons.Default.Person, contentDescription = null) },
                        label = {
                            Text(state.language.t("Profil", "Profile", "الملف الشخصي"))
                        },
                    )
                }
            }
        },
    ) { padding ->
        Box(Modifier.padding(padding)) {
            when (state.page) {
                AppPage.Home -> HomeScreen(state, viewModel)
                AppPage.Chat -> ChatScreen(state, viewModel)
                AppPage.Profile -> ProfileScreen(state, viewModel, activity)
                AppPage.Plans -> PlansScreen(state, viewModel, activity)
            }
        }
    }
}

@Composable
private fun HomeScreen(state: AppState, viewModel: KuranAIViewModel) {
    val l = state.language
    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .statusBarsPadding(),
        contentPadding = PaddingValues(22.dp),
        verticalArrangement = Arrangement.spacedBy(18.dp),
    ) {
        item {
            Row(verticalAlignment = Alignment.CenterVertically, modifier = Modifier.fillMaxWidth()) {
                QuranMark(size = 48)
                Spacer(Modifier.width(12.dp))
                Column(Modifier.weight(1f)) {
                    Text(l.t("Kuran AI", "Quran AI", "قرآن AI"), fontWeight = FontWeight.Bold, fontSize = 20.sp)
                    Text(
                        l.t("Vahiy ile konuş.", "Speak with revelation.", "تحدّث مع الوحي."),
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
                }
                LanguageMenu(l, viewModel::setLanguage)
            }
        }
        item {
            Card(
                onClick = { viewModel.newChat() },
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant),
                shape = RoundedCornerShape(28.dp),
            ) {
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(28.dp),
                    horizontalAlignment = Alignment.CenterHorizontally,
                ) {
                    FilledIconButton(
                        onClick = { viewModel.newChat() },
                        modifier = Modifier.size(74.dp),
                    ) {
                        Icon(Icons.Default.Mic, contentDescription = null, modifier = Modifier.size(34.dp))
                    }
                    Spacer(Modifier.height(14.dp))
                    Text(
                        l.t(
                            "Bugün neyi keşfetmek istersiniz?",
                            "What would you like to explore today?",
                            "ماذا تريد أن تكتشف اليوم؟",
                        ),
                        textAlign = TextAlign.Center,
                        fontWeight = FontWeight.Bold,
                        fontSize = 24.sp,
                    )
                    Text(
                        l.t(
                            "Sesli sorun veya yazarak başlayın",
                            "Ask by voice or start typing",
                            "اسأل بصوتك أو ابدأ بالكتابة",
                        ),
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
                }
            }
        }
        item {
            Text(l.t("Kaynak", "Source", "المصدر"), fontWeight = FontWeight.Bold)
            Spacer(Modifier.height(8.dp))
            LazyRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                items(SourceType.entries) { source ->
                    AssistChip(
                        onClick = { viewModel.setSource(source) },
                        label = {
                            Text(
                                when (source) {
                                    SourceType.Quran -> l.t("Kur'an Meali", "Quran Translation", "ترجمة القرآن")
                                    SourceType.Tafsir -> l.t("Tefsir", "Tafsir", "التفسير")
                                    SourceType.Hadith -> l.t("Hadis", "Hadith", "الحديث")
                                },
                            )
                        },
                        leadingIcon = { Text(source.mark, color = KuranGold, fontWeight = FontWeight.Bold) },
                        border = if (state.source == source) {
                            ButtonDefaults.outlinedButtonBorder(enabled = true)
                        } else {
                            null
                        },
                    )
                }
            }
        }
        state.dailyVerse?.let { verse ->
            item {
                SectionCard(title = l.t("Günün Ayeti", "Verse of the Day", "آية اليوم")) {
                    Text(
                        verse.arabicText,
                        fontSize = 21.sp,
                        textAlign = TextAlign.End,
                        modifier = Modifier.fillMaxWidth(),
                    )
                    Spacer(Modifier.height(10.dp))
                    Text(verse.translation, fontStyle = FontStyle.Italic)
                    Text(
                        verse.reference,
                        color = KuranGold,
                        modifier = Modifier.align(Alignment.End),
                    )
                }
            }
        }
        item {
            SectionCard(
                title = l.t("Namaz Vakitleri", "Prayer Times", "مواقيت الصلاة") +
                    state.prayerPlace.takeIf(String::isNotBlank)?.let { " · $it" }.orEmpty(),
            ) {
                if (state.prayerTimes.isEmpty()) {
                    CircularProgressIndicator(Modifier.size(22.dp))
                } else {
                    state.prayerTimes.forEach { prayer ->
                        Row(Modifier.fillMaxWidth()) {
                            Text(prayer.name, color = MaterialTheme.colorScheme.onSurfaceVariant)
                            Spacer(Modifier.weight(1f))
                            Text(prayer.time, fontWeight = FontWeight.SemiBold)
                        }
                    }
                }
            }
        }
        item {
            SectionCard(title = l.t("Hızlı Sorular", "Quick Questions", "أسئلة سريعة")) {
                quickQuestions(l).forEach { question ->
                    OutlinedButton(
                        onClick = { viewModel.newChat(question) },
                        modifier = Modifier.fillMaxWidth(),
                    ) {
                        Text(question, modifier = Modifier.weight(1f), textAlign = TextAlign.Start)
                        Icon(Icons.AutoMirrored.Filled.Send, contentDescription = null)
                    }
                }
            }
        }
    }
}

@Composable
private fun SectionCard(
    title: String,
    content: @Composable ColumnScope.() -> Unit,
) {
    OutlinedCard(shape = RoundedCornerShape(20.dp)) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(18.dp),
            verticalArrangement = Arrangement.spacedBy(7.dp),
        ) {
            Text(title, color = KuranGold, fontWeight = FontWeight.Bold)
            content()
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun ChatScreen(state: AppState, viewModel: KuranAIViewModel) {
    val l = state.language
    val drawer = rememberDrawerState(DrawerValue.Closed)
    val scope = rememberCoroutineScope()
    var input by remember { mutableStateOf("") }
    var consentAction by remember { mutableStateOf<String?>(null) }
    val listState = rememberLazyListState()
    val microphonePermission = rememberLauncherForActivityResult(
        ActivityResultContracts.RequestPermission(),
    ) { granted ->
        if (granted) {
            if (viewModel.hasAiConsent) viewModel.toggleRecording() else consentAction = "record"
        }
    }

    LaunchedEffect(state.messages.size, state.messages.lastOrNull()?.content) {
        if (state.messages.isNotEmpty()) listState.animateScrollToItem(state.messages.lastIndex)
    }

    ModalNavigationDrawer(
        drawerState = drawer,
        drawerContent = {
            Surface(
                modifier = Modifier
                    .fillMaxHeight()
                    .fillMaxWidth(0.86f),
                color = MaterialTheme.colorScheme.surface,
            ) {
                LazyColumn(
                    contentPadding = PaddingValues(top = 48.dp, bottom = 24.dp),
                ) {
                    item {
                        ListItem(
                            headlineContent = {
                                Text(l.t("Sohbet geçmişi", "Chat history", "سجل المحادثات"))
                            },
                            trailingContent = {
                                IconButton(onClick = { scope.launch { drawer.close() } }) {
                                    Icon(Icons.Default.Close, contentDescription = null)
                                }
                            },
                        )
                        ListItem(
                            headlineContent = {
                                Text(l.t("Yeni sohbet", "New chat", "محادثة جديدة"))
                            },
                            leadingContent = { Icon(Icons.Default.Add, contentDescription = null) },
                            modifier = Modifier.clickable {
                                viewModel.newChat()
                                scope.launch { drawer.close() }
                            },
                        )
                        HorizontalDivider()
                    }
                    items(state.sessions, key = { it.id }) { session ->
                        ListItem(
                            headlineContent = {
                                Text(
                                    session.firstMessage.ifBlank {
                                        l.t("Sohbet", "Chat", "محادثة")
                                    },
                                    maxLines = 2,
                                    overflow = TextOverflow.Ellipsis,
                                )
                            },
                            supportingContent = { Text(session.updatedAt.take(10)) },
                            modifier = Modifier.clickable {
                                viewModel.selectSession(session.id)
                                scope.launch { drawer.close() }
                            },
                        )
                    }
                }
            }
        },
    ) {
        Scaffold(
            topBar = {
                TopAppBar(
                    title = {
                        Column {
                            Text(l.t("Kuran AI", "Quran AI", "قرآن AI"), fontWeight = FontWeight.Bold)
                            Text(
                                l.t("Kaynak destekli asistan", "Source-assisted AI", "مساعد مدعوم بالمصادر"),
                                style = MaterialTheme.typography.labelSmall,
                            )
                        }
                    },
                    navigationIcon = {
                        IconButton(onClick = { scope.launch { drawer.open() } }) {
                            Icon(Icons.Default.Menu, contentDescription = l.t("Menü", "Menu", "القائمة"))
                        }
                    },
                    actions = {
                        IconButton(onClick = { viewModel.newChat() }) {
                            Icon(Icons.Default.Add, contentDescription = l.t("Yeni sohbet", "New chat", "محادثة جديدة"))
                        }
                    },
                )
            },
        ) { innerPadding ->
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(innerPadding),
            ) {
                LazyRow(
                    modifier = Modifier.padding(horizontal = 12.dp, vertical = 6.dp),
                    horizontalArrangement = Arrangement.spacedBy(7.dp),
                ) {
                    items(SourceType.entries) { source ->
                        AssistChip(
                            onClick = { viewModel.setSource(source) },
                            label = {
                                Text(
                                    when (source) {
                                        SourceType.Quran -> l.t("Kur'an", "Quran", "القرآن")
                                        SourceType.Tafsir -> l.t("Tefsir", "Tafsir", "التفسير")
                                        SourceType.Hadith -> l.t("Hadis", "Hadith", "الحديث")
                                    },
                                )
                            },
                            leadingIcon = { Text(source.mark, color = KuranGold) },
                        )
                    }
                }
                LazyColumn(
                    state = listState,
                    modifier = Modifier
                        .weight(1f)
                        .fillMaxWidth(),
                    contentPadding = PaddingValues(14.dp),
                    verticalArrangement = Arrangement.spacedBy(10.dp),
                ) {
                    if (state.messages.isEmpty()) {
                        item {
                            EmptyChat(l) { input = it }
                        }
                    }
                    items(state.messages, key = { it.id }) { message ->
                        MessageBubble(message, l) {
                            viewModel.speak(message.content)
                        }
                    }
                }
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 12.dp, vertical = 8.dp),
                    verticalAlignment = Alignment.Bottom,
                    horizontalArrangement = Arrangement.spacedBy(8.dp),
                ) {
                    OutlinedTextField(
                        value = input,
                        onValueChange = { input = it.take(2000) },
                        placeholder = {
                            Text(
                                l.t(
                                    "Ayet, konu veya kavram hakkında sorun…",
                                    "Ask about a verse, topic, or concept…",
                                    "اسأل عن آية أو موضوع أو مفهوم…",
                                ),
                            )
                        },
                        modifier = Modifier.weight(1f),
                        maxLines = 5,
                        keyboardOptions = KeyboardOptions(imeAction = ImeAction.Send),
                        keyboardActions = KeyboardActions(
                            onSend = {
                                if (input.isNotBlank()) {
                                    if (viewModel.hasAiConsent) {
                                        viewModel.sendMessage(input)
                                        input = ""
                                    } else {
                                        consentAction = "send"
                                    }
                                }
                            },
                        ),
                    )
                    if (input.isBlank()) {
                        FilledIconButton(
                            onClick = {
                                microphonePermission.launch(Manifest.permission.RECORD_AUDIO)
                            },
                            colors = IconButtonDefaults.filledIconButtonColors(
                                containerColor = if (state.recording) Color.Red else KuranGold,
                            ),
                            modifier = Modifier.size(52.dp),
                        ) {
                            Icon(
                                if (state.recording) Icons.Default.StopCircle else Icons.Default.Mic,
                                contentDescription = l.t("Sesli soru", "Voice question", "سؤال صوتي"),
                            )
                        }
                    } else {
                        FilledIconButton(
                            onClick = {
                                if (viewModel.hasAiConsent) {
                                    viewModel.sendMessage(input)
                                    input = ""
                                } else {
                                    consentAction = "send"
                                }
                            },
                            enabled = !state.sending,
                            modifier = Modifier.size(52.dp),
                        ) {
                            Icon(Icons.AutoMirrored.Filled.Send, contentDescription = l.t("Gönder", "Send", "إرسال"))
                        }
                    }
                }
                Text(
                    l.t(
                        "Yanıtlar eğitim amaçlıdır; fetva niteliği taşımaz.",
                        "Responses are educational and are not religious rulings.",
                        "الإجابات تعليمية وليست فتاوى شرعية.",
                    ),
                    style = MaterialTheme.typography.labelSmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant,
                    textAlign = TextAlign.Center,
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(bottom = 6.dp),
                )
            }
        }
    }

    if (consentAction != null) {
        AlertDialog(
            onDismissRequest = { consentAction = null },
            title = {
                Text(l.t("AI veri paylaşımı", "AI data sharing", "مشاركة البيانات مع الذكاء الاصطناعي"))
            },
            text = {
                Text(
                    l.t(
                        "Mesaj metniniz ve sesli soru kullandığınızda ses kaydınız, yanıt üretmek veya sesi yazıya çevirmek için OpenAI, Inc. ile paylaşılır. Adınız, e-posta adresiniz ve hesap bilgileriniz gönderilmez.",
                        "Your message text and voice recording are shared with OpenAI, Inc. to generate a response or transcribe audio. Your name, email, and account details are not sent.",
                        "تتم مشاركة نص رسالتك وتسجيلك الصوتي مع OpenAI, Inc. لإنشاء الرد أو تحويل الصوت إلى نص. لا يتم إرسال اسمك أو بريدك أو بيانات حسابك.",
                    ),
                )
            },
            confirmButton = {
                Button(onClick = {
                    viewModel.grantAiConsent()
                    if (consentAction == "send") {
                        viewModel.sendMessage(input)
                        input = ""
                    } else {
                        viewModel.toggleRecording()
                    }
                    consentAction = null
                }) {
                    Text(l.t("İzin ver ve devam et", "Allow and continue", "السماح والمتابعة"))
                }
            },
            dismissButton = {
                TextButton(onClick = { consentAction = null }) {
                    Text(l.t("Şimdi değil", "Not now", "ليس الآن"))
                }
            },
        )
    }
}

@Composable
private fun EmptyChat(language: AppLanguage, onSuggestion: (String) -> Unit) {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 42.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(14.dp),
    ) {
        QuranMark(size = 64)
        Text(
            language.t(
                "Bugün neyi keşfetmek istersiniz?",
                "What would you like to explore today?",
                "ماذا تريد أن تكتشف اليوم؟",
            ),
            textAlign = TextAlign.Center,
            fontWeight = FontWeight.Bold,
            fontSize = 24.sp,
        )
        quickQuestions(language).take(2).forEach { suggestion ->
            OutlinedButton(onClick = { onSuggestion(suggestion) }) {
                Text(suggestion)
            }
        }
    }
}

@Composable
private fun MessageBubble(
    message: ChatMessage,
    language: AppLanguage,
    onSpeak: () -> Unit,
) {
    val user = message.role == "user"
    Column(
        modifier = Modifier.fillMaxWidth(),
        horizontalAlignment = if (user) Alignment.End else Alignment.Start,
    ) {
        Text(
            if (user) language.t("Siz", "You", "أنت") else "Kuran AI",
            style = MaterialTheme.typography.labelSmall,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
        )
        Card(
            colors = CardDefaults.cardColors(
                containerColor = if (user) KuranGold else MaterialTheme.colorScheme.surfaceVariant,
            ),
            modifier = Modifier.fillMaxWidth(0.88f),
        ) {
            Column(Modifier.padding(14.dp)) {
                Text(
                    message.content.ifBlank {
                        language.t("Düşünüyor…", "Thinking…", "يفكّر…")
                    },
                    color = if (user) Color(0xFF171105) else MaterialTheme.colorScheme.onSurface,
                )
                if (!user && !message.streaming && message.content.isNotBlank()) {
                    TextButton(onClick = onSpeak, modifier = Modifier.align(Alignment.End)) {
                        Icon(Icons.Default.PlayArrow, contentDescription = null)
                        Text(language.t("Sesli oku", "Read aloud", "قراءة بصوت عالٍ"))
                    }
                }
            }
        }
    }
}

@Composable
private fun ProfileScreen(
    state: AppState,
    viewModel: KuranAIViewModel,
    activity: Activity,
) {
    val l = state.language
    var username by remember(state.user?.username) { mutableStateOf(state.user?.username.orEmpty()) }
    var currentPassword by remember { mutableStateOf("") }
    var newPassword by remember { mutableStateOf("") }
    var confirmPassword by remember { mutableStateOf("") }
    var confirmDelete by remember { mutableStateOf(false) }
    val imagePicker = rememberLauncherForActivityResult(
        ActivityResultContracts.PickVisualMedia(),
    ) { uri ->
        if (uri != null) viewModel.uploadAvatar(uri)
    }

    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .statusBarsPadding(),
        contentPadding = PaddingValues(20.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp),
    ) {
        item {
            Row(verticalAlignment = Alignment.CenterVertically, modifier = Modifier.fillMaxWidth()) {
                Column(Modifier.weight(1f)) {
                    Text(l.t("Profil", "Profile", "الملف الشخصي"), fontSize = 28.sp, fontWeight = FontWeight.Bold)
                    Text(state.user?.email.orEmpty(), color = MaterialTheme.colorScheme.onSurfaceVariant)
                }
                IconButton(onClick = { viewModel.logout(activity) }) {
                    Icon(Icons.AutoMirrored.Filled.ExitToApp, contentDescription = l.t("Çıkış", "Sign out", "تسجيل الخروج"))
                }
            }
        }
        item {
            SectionCard(l.t("Profil fotoğrafı", "Profile photo", "صورة الملف الشخصي")) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    val avatarUrl = state.user?.avatarUrl
                    if (avatarUrl != null) {
                        AsyncImage(
                            model = avatarUrl,
                            contentDescription = l.t("Profil fotoğrafı", "Profile photo", "صورة الملف الشخصي"),
                            modifier = Modifier
                                .size(76.dp)
                                .clip(CircleShape),
                        )
                    } else {
                        Icon(
                            Icons.Default.AccountCircle,
                            contentDescription = null,
                            modifier = Modifier.size(76.dp),
                        )
                    }
                    Spacer(Modifier.width(14.dp))
                    Column {
                        Button(onClick = {
                            imagePicker.launch(
                                PickVisualMediaRequest(ActivityResultContracts.PickVisualMedia.ImageOnly),
                            )
                        }) {
                            Icon(Icons.Default.PhotoCamera, contentDescription = null)
                            Text(l.t("Fotoğraf seç", "Choose photo", "اختيار صورة"))
                        }
                        if (state.user?.avatarUrl != null) {
                            TextButton(onClick = viewModel::deleteAvatar) {
                                Text(l.t("Kaldır", "Remove", "إزالة"))
                            }
                        }
                    }
                }
            }
        }
        item {
            SectionCard(l.t("Kullanıcı adı", "Username", "اسم المستخدم")) {
                OutlinedTextField(
                    value = username,
                    onValueChange = { username = it },
                    modifier = Modifier.fillMaxWidth(),
                    singleLine = true,
                )
                Button(
                    onClick = { viewModel.updateUsername(username) },
                    enabled = username.length >= 2,
                    modifier = Modifier.align(Alignment.End),
                ) {
                    Text(l.t("Kaydet", "Save", "حفظ"))
                }
            }
        }
        item {
            SectionCard(l.t("Geçiş anahtarı", "Passkey", "مفتاح المرور")) {
                Text(
                    if (state.passkeyCount > 0) {
                        l.t(
                            "Parolasız giriş etkin (${state.passkeyCount})",
                            "Passwordless sign-in is active (${state.passkeyCount})",
                            "تسجيل الدخول دون كلمة مرور مفعّل (${state.passkeyCount})",
                        )
                    } else {
                        l.t(
                            "Parmak izi veya ekran kilidi ile güvenli giriş",
                            "Secure sign-in with biometrics or screen lock",
                            "دخول آمن بالبصمة أو قفل الشاشة",
                        )
                    },
                )
                Button(onClick = { viewModel.createPasskey(activity) }) {
                    Icon(Icons.Default.Fingerprint, contentDescription = null)
                    Text(
                        if (state.passkeyCount > 0) {
                            l.t("Yeni geçiş anahtarı ekle", "Add another passkey", "إضافة مفتاح مرور آخر")
                        } else {
                            l.t("Geçiş anahtarı oluştur", "Create a passkey", "إنشاء مفتاح مرور")
                        },
                    )
                }
                if (state.passkeyCount > 0) {
                    TextButton(onClick = viewModel::deletePasskeys) {
                        Text(l.t("Tümünü kaldır", "Remove all", "إزالة الكل"))
                    }
                }
            }
        }
        item {
            SectionCard(l.t("Şifre", "Password", "كلمة المرور")) {
                if (state.user?.hasPassword == true) {
                    PasswordField(
                        currentPassword,
                        { currentPassword = it },
                        l.t("Mevcut şifre", "Current password", "كلمة المرور الحالية"),
                    )
                }
                PasswordField(
                    newPassword,
                    { newPassword = it },
                    l.t("Yeni şifre", "New password", "كلمة المرور الجديدة"),
                )
                PasswordField(
                    confirmPassword,
                    { confirmPassword = it },
                    l.t("Yeni şifre tekrar", "Confirm password", "تأكيد كلمة المرور"),
                )
                Button(
                    onClick = {
                        viewModel.updatePassword(currentPassword, newPassword)
                        currentPassword = ""
                        newPassword = ""
                        confirmPassword = ""
                    },
                    enabled = newPassword.length >= 6 && newPassword == confirmPassword,
                    modifier = Modifier.align(Alignment.End),
                ) {
                    Text(l.t("Şifreyi güncelle", "Update password", "تحديث كلمة المرور"))
                }
            }
        }
        item {
            SectionCard(l.t("Görünüm ve dil", "Appearance and language", "المظهر واللغة")) {
                Row(verticalAlignment = Alignment.CenterVertically, modifier = Modifier.fillMaxWidth()) {
                    Icon(
                        if (state.darkTheme) Icons.Default.DarkMode else Icons.Default.LightMode,
                        contentDescription = null,
                    )
                    Spacer(Modifier.width(10.dp))
                    Text(
                        l.t("Koyu tema", "Dark theme", "الوضع الداكن"),
                        modifier = Modifier.weight(1f),
                    )
                    Switch(checked = state.darkTheme, onCheckedChange = viewModel::setDarkTheme)
                }
                LanguageMenu(l, viewModel::setLanguage)
            }
        }
        item {
            SectionCard(l.t("Abonelik", "Subscription", "الاشتراك")) {
                Text(
                    state.subscription.currentPackage?.localizedName(l)
                        ?: l.t("Ücretsiz paket", "Free plan", "الباقة المجانية"),
                    fontSize = 20.sp,
                    fontWeight = FontWeight.Bold,
                )
                Text(
                    "${state.subscription.messageUsed} / ${limitText(state.subscription.messageLimit)} " +
                        l.t("soru bugün", "questions today", "سؤال اليوم"),
                    color = MaterialTheme.colorScheme.onSurfaceVariant,
                )
                Button(onClick = { viewModel.setPage(AppPage.Plans) }) {
                    Text(l.t("Paketleri görüntüle", "View plans", "عرض الباقات"))
                }
            }
        }
        item {
            SectionCard(l.t("Yasal ve destek", "Legal and support", "القانون والدعم")) {
                LinkButton(l.t("Gizlilik Politikası", "Privacy Policy", "سياسة الخصوصية"), activity, "https://kuranai.talhacan.com/privacy-policy")
                LinkButton(l.t("Kullanım Koşulları", "Terms of Use", "شروط الاستخدام"), activity, "https://kuranai.talhacan.com/terms")
                LinkButton(l.t("Destek", "Support", "الدعم"), activity, "https://kuranai.talhacan.com/support")
            }
        }
        item {
            OutlinedButton(
                onClick = { confirmDelete = true },
                colors = ButtonDefaults.outlinedButtonColors(contentColor = MaterialTheme.colorScheme.error),
                modifier = Modifier.fillMaxWidth(),
            ) {
                Icon(Icons.Default.DeleteForever, contentDescription = null)
                Text(l.t("Hesabımı kalıcı olarak sil", "Permanently delete my account", "حذف حسابي نهائياً"))
            }
        }
    }

    if (confirmDelete) {
        AlertDialog(
            onDismissRequest = { confirmDelete = false },
            title = { Text(l.t("Hesabı sil", "Delete account", "حذف الحساب")) },
            text = {
                Text(
                    l.t(
                        "Bu işlem geri alınamaz. Hesabınız ve tüm verileriniz kalıcı olarak silinecektir.",
                        "This cannot be undone. Your account and all data will be permanently deleted.",
                        "لا يمكن التراجع عن هذا الإجراء. سيتم حذف حسابك وجميع بياناتك نهائياً.",
                    ),
                )
            },
            confirmButton = {
                Button(
                    onClick = {
                        confirmDelete = false
                        viewModel.deleteAccount()
                    },
                    colors = ButtonDefaults.buttonColors(containerColor = MaterialTheme.colorScheme.error),
                ) {
                    Text(l.t("Sil", "Delete", "حذف"))
                }
            },
            dismissButton = {
                TextButton(onClick = { confirmDelete = false }) {
                    Text(l.t("İptal", "Cancel", "إلغاء"))
                }
            },
        )
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun PlansScreen(
    state: AppState,
    viewModel: KuranAIViewModel,
    activity: Activity,
) {
    val l = state.language
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(l.t("Sınırsız Erişim", "Unlimited Access", "وصول بلا حدود")) },
                navigationIcon = {
                    IconButton(onClick = { viewModel.setPage(AppPage.Profile) }) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = null)
                    }
                },
            )
        },
    ) { padding ->
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding),
            contentPadding = PaddingValues(18.dp),
            verticalArrangement = Arrangement.spacedBy(14.dp),
        ) {
            item {
                Text(
                    l.t(
                        "Paket seçin ve özelliklerin kilidini açın",
                        "Choose a plan and unlock all features",
                        "اختر باقة وافتح جميع الميزات",
                    ),
                    color = MaterialTheme.colorScheme.onSurfaceVariant,
                )
            }
            items(state.packages, key = { it.id }) { subscriptionPackage ->
                PlanCard(state, subscriptionPackage, activity, viewModel.billingStore)
            }
            item {
                OutlinedButton(
                    onClick = viewModel.billingStore::restore,
                    modifier = Modifier.fillMaxWidth(),
                ) {
                    Icon(Icons.Default.Restore, contentDescription = null)
                    Text(l.t("Satın Almaları Geri Yükle", "Restore Purchases", "استعادة المشتريات"))
                }
            }
            item {
                Text(
                    l.t(
                        "Abonelikler Google Play tarafından otomatik yenilenir. Google Play abonelik ayarlarından istediğiniz zaman iptal edebilirsiniz.",
                        "Subscriptions renew automatically through Google Play. You can cancel any time in Google Play subscription settings.",
                        "تتجدد الاشتراكات تلقائياً عبر Google Play. يمكنك الإلغاء في أي وقت من إعدادات اشتراكات Google Play.",
                    ),
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant,
                    textAlign = TextAlign.Center,
                )
                LegalLinks(l, activity)
            }
        }
    }
}

@Composable
private fun PlanCard(
    state: AppState,
    subscriptionPackage: SubscriptionPackage,
    activity: Activity,
    billing: BillingStore,
) {
    val l = state.language
    val slug = subscriptionPackage.slug.lowercase()
    val monthlyId = "kuranai.$slug.monthly"
    val yearlyId = "kuranai.$slug.yearly"
    val monthly = state.storeProducts[monthlyId]?.formattedPrice
        ?.takeIf(String::isNotBlank)
        ?: fallbackTryPrice(subscriptionPackage.monthlyPriceTl, l)
    val yearly = state.storeProducts[yearlyId]?.formattedPrice
        ?.takeIf(String::isNotBlank)
        ?: fallbackTryPrice(subscriptionPackage.yearlyPriceTl, l)
    val current = state.subscription.currentPackage?.slug == subscriptionPackage.slug

    OutlinedCard(
        colors = CardDefaults.outlinedCardColors(
            containerColor = if (slug == "diamond") {
                MaterialTheme.colorScheme.surfaceVariant
            } else {
                MaterialTheme.colorScheme.surface
            },
        ),
        border = BorderStroke(
            width = if (current) 2.dp else 1.dp,
            color = if (current) KuranGold else MaterialTheme.colorScheme.outlineVariant,
        ),
        shape = RoundedCornerShape(22.dp),
    ) {
        Column(
            modifier = Modifier.padding(18.dp),
            verticalArrangement = Arrangement.spacedBy(10.dp),
        ) {
            Row(Modifier.fillMaxWidth(), verticalAlignment = Alignment.CenterVertically) {
                Text(
                    subscriptionPackage.localizedName(l),
                    fontSize = 22.sp,
                    fontWeight = FontWeight.Bold,
                    modifier = Modifier.weight(1f),
                )
                if (slug == "diamond") {
                    Text(
                        l.t("Önerilen", "Recommended", "موصى بها"),
                        color = Color(0xFF171105),
                        modifier = Modifier
                            .background(KuranGold, RoundedCornerShape(8.dp))
                            .padding(horizontal = 8.dp, vertical = 4.dp),
                    )
                }
            }
            Text(
                "${subscriptionPackage.dailyQuestionLimit?.toString() ?: "∞"} " +
                    l.t("soru/gün", "questions/day", "سؤال/يوم") +
                    "  ·  ${subscriptionPackage.dailyTtsLimit?.toString() ?: "∞"} " +
                    l.t("TTS/gün", "TTS/day", "قراءة/يوم"),
                color = MaterialTheme.colorScheme.onSurfaceVariant,
            )
            PurchaseRow(
                price = monthly,
                period = l.t("ay", "month", "شهر"),
                label = l.t("Seç", "Choose", "اختيار"),
                onClick = { billing.purchase(activity, monthlyId) },
            )
            PurchaseRow(
                price = yearly,
                period = l.t("yıl", "year", "سنة"),
                label = l.t("Seç", "Choose", "اختيار"),
                emphasized = true,
                onClick = { billing.purchase(activity, yearlyId) },
            )
        }
    }
}

@Composable
private fun PurchaseRow(
    price: String,
    period: String,
    label: String,
    emphasized: Boolean = false,
    onClick: () -> Unit,
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .background(
                if (emphasized) KuranGold.copy(alpha = 0.12f) else MaterialTheme.colorScheme.surfaceVariant,
                RoundedCornerShape(14.dp),
            )
            .padding(14.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Column(Modifier.weight(1f)) {
            Text(price, fontWeight = FontWeight.Bold, fontSize = 18.sp)
            Text("/ $period", color = MaterialTheme.colorScheme.onSurfaceVariant)
        }
        Button(onClick = onClick) { Text(label) }
    }
}

@Composable
private fun PasswordField(
    value: String,
    onValueChange: (String) -> Unit,
    label: String,
) {
    OutlinedTextField(
        value = value,
        onValueChange = onValueChange,
        label = { Text(label) },
        visualTransformation = PasswordVisualTransformation(),
        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
        modifier = Modifier.fillMaxWidth(),
        singleLine = true,
    )
}

@Composable
private fun LegalLinks(language: AppLanguage, activity: Activity) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.Center,
    ) {
        TextButton(onClick = { openUrl(activity, "https://kuranai.talhacan.com/privacy-policy") }) {
            Text(language.t("Gizlilik", "Privacy", "الخصوصية"))
        }
        TextButton(onClick = { openUrl(activity, "https://kuranai.talhacan.com/terms") }) {
            Text(language.t("Koşullar", "Terms", "الشروط"))
        }
    }
}

@Composable
private fun LinkButton(label: String, activity: Activity, url: String) {
    TextButton(onClick = { openUrl(activity, url) }, modifier = Modifier.fillMaxWidth()) {
        Text(label, modifier = Modifier.weight(1f), textAlign = TextAlign.Start)
        Icon(Icons.AutoMirrored.Filled.Send, contentDescription = null)
    }
}

@Composable
private fun QuranMark(size: Int) {
    Box(
        modifier = Modifier
            .size(size.dp)
            .background(KuranGold.copy(alpha = 0.15f), CircleShape)
            .semantics { contentDescription = "Kuran AI" },
        contentAlignment = Alignment.Center,
    ) {
        Text("ق", color = KuranGold, fontSize = (size * 0.45).sp, fontWeight = FontWeight.Bold)
    }
}

private fun quickQuestions(language: AppLanguage): List<String> = listOf(
    language.t(
        "Sabır hakkında Kur'an ne söyler?",
        "What does the Quran say about patience?",
        "ماذا يقول القرآن عن الصبر؟",
    ),
    language.t(
        "Fâtiha suresinin ana mesajı nedir?",
        "What is the main message of Surah Al-Fatiha?",
        "ما الرسالة الأساسية لسورة الفاتحة؟",
    ),
    language.t(
        "Niyet ile ilgili sahih hadisleri açıklar mısın?",
        "Can you explain authentic hadiths about intention?",
        "هل تشرح الأحاديث الصحيحة عن النية؟",
    ),
)

private fun limitText(value: Int): String =
    if (value < 0 || value >= 999_999) "∞" else value.toString()

private fun fallbackTryPrice(value: Double, language: AppLanguage): String {
    val locale = when (language) {
        AppLanguage.Turkish -> Locale.forLanguageTag("tr-TR")
        AppLanguage.English -> Locale.US
        AppLanguage.Arabic -> Locale.forLanguageTag("ar-SA")
    }
    return NumberFormat.getCurrencyInstance(locale).apply {
        currency = Currency.getInstance("TRY")
    }.format(value)
}

private fun openUrl(activity: Activity, url: String) {
    activity.startActivity(Intent(Intent.ACTION_VIEW, Uri.parse(url)))
}
