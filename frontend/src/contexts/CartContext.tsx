import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { CartItem, CartContextType, TicketTier } from '../types'

const CartContext = createContext<CartContextType | undefined>(undefined)

export const useCart = (): CartContextType => {
  const context = useContext(CartContext)
  if (!context) {
    throw new Error('useCart must be used within a CartProvider')
  }
  return context
}

interface CartProviderProps {
  children: ReactNode
}

export const CartProvider: React.FC<CartProviderProps> = ({ children }) => {
  const [items, setItems] = useState<CartItem[]>([])
  const [isOpen, setIsOpen] = useState<boolean>(false)

  // Load cart from localStorage on mount
  useEffect(() => {
    const savedCart = localStorage.getItem('cart')
    if (savedCart) {
      try {
        setItems(JSON.parse(savedCart))
      } catch (error) {
        console.error('Failed to load cart from localStorage:', error)
      }
    }
  }, [])

  // Save cart to localStorage whenever items change
  useEffect(() => {
    localStorage.setItem('cart', JSON.stringify(items))
  }, [items])

  const addItem = (eventId: string, ticketTier: TicketTier, quantity: number = 1): void => {
    setItems(prevItems => {
      const existingItem = prevItems.find(
        item => item.eventId === eventId && item.ticketTier.id === ticketTier.id
      )

      if (existingItem) {
        // Update quantity if item already exists
        return prevItems.map(item =>
          item.eventId === eventId && item.ticketTier.id === ticketTier.id
            ? { ...item, quantity: item.quantity + quantity }
            : item
        )
      } else {
        // Add new item
        return [...prevItems, {
          id: `${eventId}-${ticketTier.id}`,
          eventId,
          ticketTier,
          quantity,
          addedAt: new Date().toISOString()
        }]
      }
    })
  }

  const removeItem = (itemId: string): void => {
    setItems(prevItems => prevItems.filter(item => item.id !== itemId))
  }

  const updateQuantity = (itemId: string, quantity: number): void => {
    if (quantity <= 0) {
      removeItem(itemId)
      return
    }

    setItems(prevItems =>
      prevItems.map(item =>
        item.id === itemId ? { ...item, quantity } : item
      )
    )
  }

  const clearCart = (): void => {
    setItems([])
  }

  const getTotalItems = (): number => {
    return items.reduce((total, item) => total + item.quantity, 0)
  }

  const getTotalPrice = (): number => {
    return items.reduce((total, item) => {
      return total + (item.ticketTier.price * item.quantity)
    }, 0)
  }

  const getItemsByEvent = (eventId: string): CartItem[] => {
    return items.filter(item => item.eventId === eventId)
  }

  const hasItems = (): boolean => {
    return items.length > 0
  }

  const openCart = (): void => setIsOpen(true)
  const closeCart = (): void => setIsOpen(false)
  const toggleCart = (): void => setIsOpen(!isOpen)

  const value: CartContextType = {
    items,
    isOpen,
    addItem,
    removeItem,
    updateQuantity,
    clearCart,
    getTotalItems,
    getTotalPrice,
    getItemsByEvent,
    hasItems,
    openCart,
    closeCart,
    toggleCart
  }

  return (
    <CartContext.Provider value={value}>
      {children}
    </CartContext.Provider>
  )
}

