import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import Landing from './pages/Landing'
import Privacy from './pages/Privacy'
import Terms from './pages/Terms'
import Support from './pages/Support'
import Contact from './pages/Contact'
import Chat from './pages/Chat'
import Packages from './pages/Packages'
import AppleCallback from './pages/AppleCallback'
import AndroidAppleSignIn from './pages/AndroidAppleSignIn'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/auth/apple/android" element={<AndroidAppleSignIn />} />
        <Route element={<Layout />}>
          <Route path="/" element={<Landing />} />
          <Route path="/app" element={<Chat />} />
          <Route path="/packages" element={<Packages />} />
          <Route path="/privacy-policy" element={<Privacy />} />
          <Route path="/terms" element={<Terms />} />
          <Route path="/support" element={<Support />} />
          <Route path="/contact" element={<Contact />} />
          <Route path="/auth/apple/callback" element={<AppleCallback />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
