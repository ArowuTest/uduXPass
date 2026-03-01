/*
 * GalleryUploader — uduXPass Admin Component
 * Design: Dark navy/amber brand system
 *
 * Manages an ordered array of gallery image URLs.
 * Each image can be uploaded via file or entered as a URL.
 * Supports up to 8 images with drag-to-reorder (planned).
 */

import React, { useState, useRef } from 'react'
import { Plus, X, Upload, Link, Loader2, GripVertical } from 'lucide-react'

interface GalleryUploaderProps {
  images: string[]
  onChange: (images: string[]) => void
  maxImages?: number
}

type ItemMode = 'upload' | 'url'

interface GalleryItem {
  url: string
  uploading: boolean
  error: string
  mode: ItemMode
  urlInput: string
}

const GalleryUploader: React.FC<GalleryUploaderProps> = ({
  images,
  onChange,
  maxImages = 8
}) => {
  const [items, setItems] = useState<GalleryItem[]>(
    images.map(url => ({ url, uploading: false, error: '', mode: 'upload' as ItemMode, urlInput: url }))
  )
  const fileInputRefs = useRef<(HTMLInputElement | null)[]>([])

  const syncToParent = (newItems: GalleryItem[]) => {
    onChange(newItems.map(i => i.url).filter(Boolean))
  }

  const addItem = () => {
    if (items.length >= maxImages) return
    const newItems = [...items, { url: '', uploading: false, error: '', mode: 'upload' as ItemMode, urlInput: '' }]
    setItems(newItems)
  }

  const removeItem = (index: number) => {
    const newItems = items.filter((_, i) => i !== index)
    setItems(newItems)
    syncToParent(newItems)
  }

  const uploadFile = async (index: number, file: File) => {
    const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
    if (!allowedTypes.includes(file.type)) {
      const newItems = [...items]
      newItems[index] = { ...newItems[index], error: 'Invalid file type', uploading: false }
      setItems(newItems)
      return
    }

    if (file.size > 10 * 1024 * 1024) {
      const newItems = [...items]
      newItems[index] = { ...newItems[index], error: 'File too large (max 10MB)', uploading: false }
      setItems(newItems)
      return
    }

    const newItems = [...items]
    newItems[index] = { ...newItems[index], uploading: true, error: '' }
    setItems(newItems)

    try {
      const adminToken = localStorage.getItem('adminToken')
      const formData = new FormData()
      formData.append('file', file)
      formData.append('category', 'gallery')

      const response = await fetch('/v1/admin/upload', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${adminToken}` },
        body: formData
      })

      const data = await response.json()

      if (response.ok && data.success) {
        const updated = [...items]
        updated[index] = { ...updated[index], url: data.url, uploading: false, error: '', urlInput: data.url }
        setItems(updated)
        syncToParent(updated)
      } else {
        const updated = [...items]
        updated[index] = { ...updated[index], error: data.error || 'Upload failed', uploading: false }
        setItems(updated)
      }
    } catch {
      const updated = [...items]
      updated[index] = { ...updated[index], error: 'Network error', uploading: false }
      setItems(updated)
    }
  }

  const applyUrl = (index: number) => {
    const url = items[index].urlInput.trim()
    if (!url) return
    const updated = [...items]
    updated[index] = { ...updated[index], url, error: '' }
    setItems(updated)
    syncToParent(updated)
  }

  const setMode = (index: number, mode: ItemMode) => {
    const updated = [...items]
    updated[index] = { ...updated[index], mode }
    setItems(updated)
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <label className="text-sm font-medium" style={{ color: 'var(--text-primary)' }}>
          Gallery Images
        </label>
        <span className="text-xs" style={{ color: 'var(--text-secondary)' }}>
          {items.length}/{maxImages} images
        </span>
      </div>

      {/* Gallery Grid */}
      <div className="grid grid-cols-2 gap-3">
        {items.map((item, index) => (
          <div
            key={index}
            className="relative rounded-xl border overflow-hidden"
            style={{ borderColor: 'var(--border-color)', background: 'var(--bg-card)' }}
          >
            {/* Mode toggle */}
            <div className="flex border-b" style={{ borderColor: 'var(--border-color)' }}>
              <button
                type="button"
                onClick={() => setMode(index, 'upload')}
                className="flex-1 flex items-center justify-center gap-1 py-1.5 text-xs font-medium transition-colors"
                style={{
                  background: item.mode === 'upload' ? 'rgba(245,158,11,0.15)' : 'transparent',
                  color: item.mode === 'upload' ? 'var(--accent-amber)' : 'var(--text-secondary)'
                }}
              >
                <Upload size={10} /> Upload
              </button>
              <button
                type="button"
                onClick={() => setMode(index, 'url')}
                className="flex-1 flex items-center justify-center gap-1 py-1.5 text-xs font-medium transition-colors"
                style={{
                  background: item.mode === 'url' ? 'rgba(245,158,11,0.15)' : 'transparent',
                  color: item.mode === 'url' ? 'var(--accent-amber)' : 'var(--text-secondary)'
                }}
              >
                <Link size={10} /> URL
              </button>
            </div>

            {/* Content */}
            <div className="p-2">
              {item.url ? (
                <div className="relative">
                  <img
                    src={item.url}
                    alt={`Gallery ${index + 1}`}
                    className="w-full h-28 object-cover rounded-lg"
                    onError={(e) => { (e.target as HTMLImageElement).src = 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100"><rect fill="%23374151" width="100" height="100"/></svg>' }}
                  />
                  <button
                    type="button"
                    onClick={() => removeItem(index)}
                    className="absolute top-1 right-1 rounded-full p-0.5"
                    style={{ background: 'rgba(0,0,0,0.7)' }}
                  >
                    <X size={12} className="text-white" />
                  </button>
                </div>
              ) : item.mode === 'upload' ? (
                <div>
                  <input
                    ref={el => { fileInputRefs.current[index] = el }}
                    type="file"
                    accept="image/jpeg,image/png,image/webp,image/gif"
                    onChange={e => { const f = e.target.files?.[0]; if (f) uploadFile(index, f) }}
                    className="hidden"
                  />
                  <div
                    onClick={() => !item.uploading && fileInputRefs.current[index]?.click()}
                    className="flex flex-col items-center justify-center h-28 rounded-lg border-2 border-dashed cursor-pointer transition-colors"
                    style={{ borderColor: 'var(--border-color)' }}
                  >
                    {item.uploading ? (
                      <Loader2 size={20} className="animate-spin" style={{ color: 'var(--accent-amber)' }} />
                    ) : (
                      <>
                        <Upload size={20} style={{ color: 'var(--text-secondary)' }} />
                        <span className="text-xs mt-1" style={{ color: 'var(--text-secondary)' }}>Click to upload</span>
                      </>
                    )}
                  </div>
                  {item.error && (
                    <p className="text-xs text-red-400 mt-1">{item.error}</p>
                  )}
                </div>
              ) : (
                <div className="space-y-1">
                  <input
                    type="url"
                    value={item.urlInput}
                    onChange={e => {
                      const updated = [...items]
                      updated[index] = { ...updated[index], urlInput: e.target.value }
                      setItems(updated)
                    }}
                    onKeyDown={e => e.key === 'Enter' && (e.preventDefault(), applyUrl(index))}
                    placeholder="https://..."
                    className="w-full rounded px-2 py-1.5 text-xs border"
                    style={{ background: 'var(--bg-elevated)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }}
                  />
                  <div className="flex gap-1">
                    <button
                      type="button"
                      onClick={() => applyUrl(index)}
                      className="flex-1 py-1 rounded text-xs font-medium"
                      style={{ background: 'var(--accent-amber)', color: '#0f1729' }}
                    >
                      Apply
                    </button>
                    <button
                      type="button"
                      onClick={() => removeItem(index)}
                      className="px-2 py-1 rounded text-xs"
                      style={{ background: 'var(--bg-elevated)', color: 'var(--text-secondary)' }}
                    >
                      <X size={10} />
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        ))}

        {/* Add Image Button */}
        {items.length < maxImages && (
          <button
            type="button"
            onClick={addItem}
            className="flex flex-col items-center justify-center rounded-xl border-2 border-dashed transition-all h-full min-h-[120px]"
            style={{ borderColor: 'var(--border-color)', color: 'var(--text-secondary)' }}
            onMouseEnter={e => { (e.currentTarget as HTMLElement).style.borderColor = 'var(--accent-amber)'; (e.currentTarget as HTMLElement).style.color = 'var(--accent-amber)' }}
            onMouseLeave={e => { (e.currentTarget as HTMLElement).style.borderColor = 'var(--border-color)'; (e.currentTarget as HTMLElement).style.color = 'var(--text-secondary)' }}
          >
            <Plus size={20} />
            <span className="text-xs mt-1">Add Image</span>
          </button>
        )}
      </div>

      <p className="text-xs" style={{ color: 'var(--text-secondary)' }}>
        Up to {maxImages} images · JPG, PNG, WebP, GIF · Max 10MB each
      </p>
    </div>
  )
}

export default GalleryUploader
