/*
 * MediaUploader — uduXPass Admin Component
 * Design: Dark navy/amber brand system
 *
 * Supports two input modes:
 *  1. Drag-and-drop / click-to-browse file upload → POST /v1/admin/upload
 *  2. URL paste fallback (for already-hosted assets like YouTube, Cloudinary, etc.)
 *
 * StorageProvider abstraction in the backend means this component will work
 * identically when the backend switches from LocalStorage → GCS in production.
 */

import React, { useState, useRef, useCallback } from 'react'
import { Upload, Link, X, Image as ImageIcon, Film, CheckCircle, AlertCircle, Loader2 } from 'lucide-react'

export type MediaCategory = 'event_image' | 'gallery' | 'ticket_image' | 'video'

interface MediaUploaderProps {
  label: string
  value: string
  onChange: (url: string) => void
  category?: MediaCategory
  accept?: string
  maxSizeMB?: number
  placeholder?: string
  hint?: string
  className?: string
}

type InputMode = 'upload' | 'url'
type UploadState = 'idle' | 'uploading' | 'success' | 'error'

const MediaUploader: React.FC<MediaUploaderProps> = ({
  label,
  value,
  onChange,
  category = 'event_image',
  accept = 'image/jpeg,image/png,image/webp,image/gif',
  maxSizeMB = 10,
  placeholder = 'https://example.com/image.jpg',
  hint,
  className = ''
}) => {
  const [mode, setMode] = useState<InputMode>('upload')
  const [uploadState, setUploadState] = useState<UploadState>('idle')
  const [uploadError, setUploadError] = useState<string>('')
  const [isDragging, setIsDragging] = useState(false)
  const [urlInput, setUrlInput] = useState(value || '')
  const fileInputRef = useRef<HTMLInputElement>(null)

  const isVideo = category === 'video'
  const isImage = !isVideo

  const uploadFile = useCallback(async (file: File) => {
    // Validate file type
    const allowedTypes = accept.split(',').map(t => t.trim())
    if (!allowedTypes.includes(file.type)) {
      setUploadError(`Invalid file type. Allowed: ${allowedTypes.join(', ')}`)
      setUploadState('error')
      return
    }

    // Validate file size
    const maxBytes = maxSizeMB * 1024 * 1024
    if (file.size > maxBytes) {
      setUploadError(`File too large. Maximum size: ${maxSizeMB}MB`)
      setUploadState('error')
      return
    }

    setUploadState('uploading')
    setUploadError('')

    try {
      const adminToken = localStorage.getItem('adminToken')
      const formData = new FormData()
      formData.append('file', file)
      formData.append('category', category)

      const response = await fetch('/v1/admin/upload', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${adminToken}` },
        body: formData
      })

      const data = await response.json()

      if (response.ok && data.success) {
        onChange(data.url)
        setUploadState('success')
      } else {
        setUploadError(data.error || 'Upload failed')
        setUploadState('error')
      }
    } catch (err) {
      setUploadError('Network error. Please try again.')
      setUploadState('error')
    }
  }, [accept, category, maxSizeMB, onChange])

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) uploadFile(file)
  }

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
    const file = e.dataTransfer.files?.[0]
    if (file) uploadFile(file)
  }, [uploadFile])

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(true)
  }

  const handleDragLeave = () => setIsDragging(false)

  const handleUrlApply = () => {
    if (urlInput.trim()) {
      onChange(urlInput.trim())
      setUploadState('success')
    }
  }

  const handleClear = () => {
    onChange('')
    setUrlInput('')
    setUploadState('idle')
    setUploadError('')
    if (fileInputRef.current) fileInputRef.current.value = ''
  }

  const hasValue = Boolean(value)

  return (
    <div className={`space-y-2 ${className}`}>
      {/* Label */}
      <label className="block text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
        {label}
      </label>

      {/* Mode Tabs */}
      <div className="flex rounded-lg overflow-hidden border" style={{ borderColor: 'var(--border-color)' }}>
        <button
          type="button"
          onClick={() => setMode('upload')}
          className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-sm font-medium transition-colors"
          style={{
            background: mode === 'upload' ? 'var(--accent-amber)' : 'var(--bg-card)',
            color: mode === 'upload' ? '#0f1729' : 'var(--text-secondary)'
          }}
        >
          <Upload size={14} />
          Upload File
        </button>
        <button
          type="button"
          onClick={() => setMode('url')}
          className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-sm font-medium transition-colors"
          style={{
            background: mode === 'url' ? 'var(--accent-amber)' : 'var(--bg-card)',
            color: mode === 'url' ? '#0f1729' : 'var(--text-secondary)'
          }}
        >
          <Link size={14} />
          Paste URL
        </button>
      </div>

      {/* Upload Mode */}
      {mode === 'upload' && (
        <div>
          <input
            ref={fileInputRef}
            type="file"
            accept={accept}
            onChange={handleFileChange}
            className="hidden"
          />
          <div
            onClick={() => uploadState !== 'uploading' && fileInputRef.current?.click()}
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            className="relative rounded-xl border-2 border-dashed transition-all cursor-pointer"
            style={{
              borderColor: isDragging ? 'var(--accent-amber)' : uploadState === 'error' ? '#ef4444' : 'var(--border-color)',
              background: isDragging ? 'rgba(245,158,11,0.05)' : 'var(--bg-card)',
              minHeight: '120px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}
          >
            {uploadState === 'uploading' ? (
              <div className="flex flex-col items-center gap-2 p-6">
                <Loader2 size={28} className="animate-spin" style={{ color: 'var(--accent-amber)' }} />
                <span className="text-sm" style={{ color: 'var(--text-secondary)' }}>Uploading...</span>
              </div>
            ) : uploadState === 'success' && hasValue ? (
              <div className="w-full p-3">
                {isImage ? (
                  <div className="relative">
                    <img
                      src={value}
                      alt="Uploaded"
                      className="w-full h-40 object-cover rounded-lg"
                      onError={(e) => { (e.target as HTMLImageElement).src = 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100"><rect fill="%23374151" width="100" height="100"/><text fill="%239ca3af" x="50" y="55" text-anchor="middle" font-size="12">Image</text></svg>' }}
                    />
                    <button
                      type="button"
                      onClick={(e) => { e.stopPropagation(); handleClear() }}
                      className="absolute top-2 right-2 rounded-full p-1"
                      style={{ background: 'rgba(0,0,0,0.7)' }}
                    >
                      <X size={14} className="text-white" />
                    </button>
                    <div className="mt-2 flex items-center gap-1 text-xs" style={{ color: '#22c55e' }}>
                      <CheckCircle size={12} />
                      Uploaded successfully
                    </div>
                  </div>
                ) : (
                  <div className="flex items-center gap-3 p-3">
                    <Film size={24} style={{ color: 'var(--accent-amber)' }} />
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium truncate" style={{ color: 'var(--text-primary)' }}>{value}</p>
                      <p className="text-xs" style={{ color: '#22c55e' }}>Uploaded successfully</p>
                    </div>
                    <button type="button" onClick={(e) => { e.stopPropagation(); handleClear() }}>
                      <X size={16} style={{ color: 'var(--text-secondary)' }} />
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <div className="flex flex-col items-center gap-2 p-6 text-center">
                {isImage ? (
                  <ImageIcon size={28} style={{ color: 'var(--text-secondary)' }} />
                ) : (
                  <Film size={28} style={{ color: 'var(--text-secondary)' }} />
                )}
                <div>
                  <p className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
                    Drop file here or <span style={{ color: 'var(--accent-amber)' }}>browse</span>
                  </p>
                  <p className="text-xs mt-1" style={{ color: 'var(--text-secondary)' }}>
                    {accept.replace(/image\//g, '').replace(/video\//g, '').toUpperCase().replace(/,/g, ', ')} · Max {maxSizeMB}MB
                  </p>
                </div>
              </div>
            )}
          </div>

          {uploadState === 'error' && (
            <div className="flex items-center gap-1 mt-1 text-xs text-red-400">
              <AlertCircle size={12} />
              {uploadError}
            </div>
          )}
        </div>
      )}

      {/* URL Mode */}
      {mode === 'url' && (
        <div className="space-y-2">
          <div className="flex gap-2">
            <input
              type="url"
              value={urlInput}
              onChange={(e) => setUrlInput(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), handleUrlApply())}
              placeholder={placeholder}
              className="flex-1 rounded-lg px-3 py-2 text-sm border"
              style={{
                background: 'var(--bg-card)',
                borderColor: 'var(--border-color)',
                color: 'var(--text-primary)'
              }}
            />
            <button
              type="button"
              onClick={handleUrlApply}
              className="px-4 py-2 rounded-lg text-sm font-medium"
              style={{ background: 'var(--accent-amber)', color: '#0f1729' }}
            >
              Apply
            </button>
            {value && (
              <button
                type="button"
                onClick={handleClear}
                className="px-3 py-2 rounded-lg text-sm"
                style={{ background: 'var(--bg-elevated)', color: 'var(--text-secondary)' }}
              >
                <X size={14} />
              </button>
            )}
          </div>

          {/* Preview for URL mode */}
          {value && isImage && (
            <img
              src={value}
              alt="Preview"
              className="w-full h-40 object-cover rounded-lg"
              onError={(e) => { (e.target as HTMLImageElement).style.display = 'none' }}
            />
          )}
          {value && !isImage && (
            <div className="flex items-center gap-2 p-3 rounded-lg" style={{ background: 'var(--bg-elevated)' }}>
              <Film size={16} style={{ color: 'var(--accent-amber)' }} />
              <span className="text-sm truncate" style={{ color: 'var(--text-primary)' }}>{value}</span>
            </div>
          )}
        </div>
      )}

      {hint && (
        <p className="text-xs" style={{ color: 'var(--text-secondary)' }}>{hint}</p>
      )}
    </div>
  )
}

export default MediaUploader
