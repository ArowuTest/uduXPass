import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import { CartProvider } from './contexts/CartContext'
import ErrorBoundary from './components/ErrorBoundary'
import Navbar from './components/layout/Navbar'
import Footer from './components/layout/Footer'
import AdminLayout from './components/layout/AdminLayout'
import HomePage from './pages/HomePage'
import EventsPage from './pages/EventsPage'
import EventDetailsPage from './pages/EventDetailsPage';
// Fixed auth imports - now pointing to auth folder
import LoginPage from './pages/auth/LoginPage'
import RegisterPage from './pages/auth/RegisterPage'
import CheckoutPage from './pages/CheckoutPage'
import OrderConfirmationPage from './pages/OrderConfirmationPage'
import ProfilePage from './pages/ProfilePage'
import UserTicketsPage from './pages/UserTicketsPage'
// Admin imports - updated to use TypeScript versions
import AdminLoginPage from './pages/admin/AdminLoginPage'
import AdminDashboard from './pages/admin/AdminDashboard'
import AdminEventsPage from './pages/admin/AdminEventsPage'
import AdminEventCreatePage from './pages/admin/AdminEventCreatePage'
import AdminUserManagementPage from '@/pages/admin/AdminUserManagementPage'
import RegularUserManagementPage from '@/pages/admin/RegularUserManagementPage'
import AdminOrderManagementPage from './pages/admin/AdminOrderManagementPage'
import AdminScannerManagementPage from './pages/admin/AdminScannerManagementPage'
import AdminScannerUserManagementPage from './pages/admin/AdminScannerUserManagementPage'
import AdminTicketValidationPage from './pages/admin/AdminTicketValidationPage'
import AdminAnalyticsPage from './pages/admin/AdminAnalyticsPage'
import AdminSettingsPage from './pages/admin/AdminSettingsPage'
import ProtectedRoute from './components/auth/ProtectedRoute'
import AdminProtectedRoute from './components/auth/AdminProtectedRoute'
import RoleBasedRedirect from './components/auth/RoleBasedRedirect'
import LoadingSpinner from './components/ui/LoadingSpinner'
import { Toaster } from './components/ui/toaster'
import { AdminPermission } from './types'
import './App.css'

// Component to handle role-based layout rendering
const AppContent: React.FC = () => {
  const { isAdmin, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900">
        <div className="text-center">
          <div className="mb-8">
            <h1 className="text-6xl font-bold text-white mb-2">uduXPass</h1>
            <p className="text-xl text-purple-200">Premium Event Ticketing</p>
          </div>
          <LoadingSpinner size="lg" />
        </div>
      </div>
    )
  }

  return (
    <Routes>
      {/* Root redirect based on authentication status */}
      <Route path="/" element={<RoleBasedRedirect />} />
      
      {/* Public routes (accessible to everyone) */}
      <Route path="/home" element={
        <div className="min-h-screen bg-background flex flex-col">
          <Navbar />
          <main className="flex-1">
            <HomePage />
          </main>
          <Footer />
        </div>
      } />
      
      <Route path="/events" element={
        <div className="min-h-screen bg-background flex flex-col">
          <Navbar />
          <main className="flex-1">
            <EventsPage />
          </main>
          <Footer />
        </div>
      } />
      
      <Route path="/events/:id" element={
        <div className="min-h-screen bg-background flex flex-col">
          <Navbar />
          <main className="flex-1">
            <EventDetailsPage />
          </main>
          <Footer />
        </div>
      } />
      
      {/* Authentication routes */}
      <Route path="/login" element={
        <div className="min-h-screen bg-background flex flex-col">
          <Navbar />
          <main className="flex-1">
            <LoginPage />
          </main>
          <Footer />
        </div>
      } />
      
      <Route path="/register" element={
        <div className="min-h-screen bg-background flex flex-col">
          <Navbar />
          <main className="flex-1">
            <RegisterPage />
          </main>
          <Footer />
        </div>
      } />

      {/* User-only routes */}
      <Route path="/checkout" element={
        <ProtectedRoute>
          <div className="min-h-screen bg-background flex flex-col">
            <Navbar />
            <main className="flex-1">
              <CheckoutPage />
            </main>
            <Footer />
          </div>
        </ProtectedRoute>
      } />
      
      <Route path="/order-confirmation/:orderId" element={
        <ProtectedRoute>
          <div className="min-h-screen bg-background flex flex-col">
            <Navbar />
            <main className="flex-1">
              <OrderConfirmationPage />
            </main>
            <Footer />
          </div>
        </ProtectedRoute>
      } />
      
      <Route path="/profile" element={
        <ProtectedRoute>
          <div className="min-h-screen bg-background flex flex-col">
            <Navbar />
            <main className="flex-1">
              <ProfilePage />
            </main>
            <Footer />
          </div>
        </ProtectedRoute>
      } />
      
      <Route path="/tickets" element={
        <ProtectedRoute>
          <div className="min-h-screen bg-background flex flex-col">
            <Navbar />
            <main className="flex-1">
              <UserTicketsPage />
            </main>
            <Footer />
          </div>
        </ProtectedRoute>
      } />

      {/* Admin authentication */}
      <Route path="/admin/login" element={<AdminLoginPage />} />

      {/* Admin routes with AdminLayout */}
      <Route path="/admin/dashboard" element={
        <AdminProtectedRoute>
          <AdminLayout>
            <AdminDashboard />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      <Route path="/admin/events" element={
        <AdminProtectedRoute requiredPermissions={[AdminPermission.EVENTS_VIEW]}>
          <AdminLayout>
            <AdminEventsPage />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      <Route path="/admin/events/create" element={
        <AdminProtectedRoute requiredPermissions={[AdminPermission.EVENTS_CREATE]}>
          <AdminLayout>
            <AdminEventCreatePage />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      <Route path="/admin/users" element={
        <AdminProtectedRoute requiredPermissions={[AdminPermission.USERS_VIEW]}>
          <AdminLayout>
            <RegularUserManagementPage />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      <Route path="/admin/orders" element={
        <AdminProtectedRoute requiredPermissions={[AdminPermission.ORDERS_VIEW]}>
          <AdminLayout>
            <AdminOrderManagementPage />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      <Route path="/admin/scanners" element={
        <AdminProtectedRoute requiredPermissions={[AdminPermission.SCANNERS_VIEW]}>
          <AdminLayout>
            <AdminScannerManagementPage />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      <Route path="/admin/tickets" element={
        <AdminProtectedRoute requiredPermissions={[AdminPermission.TICKETS_VIEW]}>
          <AdminLayout>
            <AdminTicketValidationPage />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      <Route path="/admin/analytics" element={
        <AdminProtectedRoute requiredPermissions={[AdminPermission.ANALYTICS_VIEW]}>
          <AdminLayout>
            <AdminAnalyticsPage />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      <Route path="/admin/settings" element={
        <AdminProtectedRoute requiredPermissions={[AdminPermission.SETTINGS_UPDATE]}>
          <AdminLayout>
            <AdminSettingsPage />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      <Route path="/admin/admin-users" element={
        <AdminProtectedRoute requiredPermissions={[AdminPermission.ADMIN_CREATE]}>
          <AdminLayout>
            <AdminUserManagementPage />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      <Route path="/admin/scanner-users" element={
        <AdminProtectedRoute requiredPermissions={[AdminPermission.SCANNERS_VIEW]}>
          <AdminLayout>
            <AdminScannerUserManagementPage />
          </AdminLayout>
        </AdminProtectedRoute>
      } />

      {/* Catch-all redirect */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

const App: React.FC = () => {
  return (
    <ErrorBoundary>
      <AuthProvider>
        <CartProvider>
          <Router>
            <AppContent />
            <Toaster />
          </Router>
        </CartProvider>
      </AuthProvider>
    </ErrorBoundary>
  )
}

export default App

