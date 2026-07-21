function base64URLToBuffer(value) {
  const base64 = value.replace(/-/g, '+').replace(/_/g, '/')
  const padded = base64.padEnd(Math.ceil(base64.length / 4) * 4, '=')
  const binary = window.atob(padded)
  return Uint8Array.from(binary, character => character.charCodeAt(0)).buffer
}

function bufferToBase64URL(value) {
  if (value == null) return null
  const bytes = new Uint8Array(value)
  let binary = ''
  bytes.forEach(byte => { binary += String.fromCharCode(byte) })
  return window.btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '')
}

function creationOptionsFromJSON(options) {
  return {
    ...options,
    challenge: base64URLToBuffer(options.challenge),
    user: { ...options.user, id: base64URLToBuffer(options.user.id) },
    excludeCredentials: (options.excludeCredentials || []).map(item => ({
      ...item,
      id: base64URLToBuffer(item.id),
    })),
  }
}

function requestOptionsFromJSON(options) {
  return {
    ...options,
    challenge: base64URLToBuffer(options.challenge),
    allowCredentials: (options.allowCredentials || []).map(item => ({
      ...item,
      id: base64URLToBuffer(item.id),
    })),
  }
}

function registrationCredentialToJSON(credential) {
  return {
    id: credential.id,
    rawId: bufferToBase64URL(credential.rawId),
    type: credential.type,
    authenticatorAttachment: credential.authenticatorAttachment,
    clientExtensionResults: credential.getClientExtensionResults(),
    response: {
      attestationObject: bufferToBase64URL(credential.response.attestationObject),
      clientDataJSON: bufferToBase64URL(credential.response.clientDataJSON),
      transports: credential.response.getTransports?.() || [],
    },
  }
}

function assertionCredentialToJSON(credential) {
  return {
    id: credential.id,
    rawId: bufferToBase64URL(credential.rawId),
    type: credential.type,
    authenticatorAttachment: credential.authenticatorAttachment,
    clientExtensionResults: credential.getClientExtensionResults(),
    response: {
      authenticatorData: bufferToBase64URL(credential.response.authenticatorData),
      clientDataJSON: bufferToBase64URL(credential.response.clientDataJSON),
      signature: bufferToBase64URL(credential.response.signature),
      userHandle: bufferToBase64URL(credential.response.userHandle),
    },
  }
}

export function passkeysSupported() {
  return typeof window !== 'undefined' &&
    window.isSecureContext &&
    typeof window.PublicKeyCredential !== 'undefined' &&
    typeof navigator.credentials?.create === 'function' &&
    typeof navigator.credentials?.get === 'function'
}

export async function createPasskey(options) {
  if (!passkeysSupported()) throw new Error('Passkeys are not supported on this device or browser.')
  const credential = await navigator.credentials.create({ publicKey: creationOptionsFromJSON(options) })
  if (!credential) throw new Error('Passkey creation was cancelled.')
  return registrationCredentialToJSON(credential)
}

export async function getPasskey(options) {
  if (!passkeysSupported()) throw new Error('Passkeys are not supported on this device or browser.')
  const credential = await navigator.credentials.get({ publicKey: requestOptionsFromJSON(options) })
  if (!credential) throw new Error('Passkey sign-in was cancelled.')
  return assertionCredentialToJSON(credential)
}
