import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { useAuth } from '../contexts/AuthContext';
import { userAPI, ordersAPI, ticketsAPI } from '../services/api';
import { Order, Ticket, UpdateProfileData, ChangePasswordData } from '../types/api';
import { formatCurrency, formatDateTime } from '../lib/utils';
import { 
  User, 
  Mail, 
  Phone, 
  Calendar, 
  MapPin, 
  Ticket as TicketIcon,
  Settings,
  Eye,
  EyeOff,
  CheckCircle,
  XCircle
} from 'lucide-react';
import { motion } from 'framer-motion';

interface ProfileFormData extends UpdateProfileData {
  email: string;
}

interface PasswordFormData extends ChangePasswordData {}

const ProfilePage: React.FC = () => {
  const { user, updateUser } = useAuth();
  const [activeTab, setActiveTab] = useState<'profile' | 'orders' | 'tickets' | 'settings'>('profile');
  const [orders, setOrders] = useState<Order[]>([]);
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isUpdating, setIsUpdating] = useState<boolean>(false);
  const [showPassword, setShowPassword] = useState<boolean>(false);
  
  const [profileForm, setProfileForm] = useState<ProfileFormData>({
    first_name: user?.first_name || '',
    last_name: user?.last_name || '',
    phone: user?.phone || '',
    email: user?.email || ''
  });
  
  const [passwordForm, setPasswordForm] = useState<PasswordFormData>({
    current_password: '',
    new_password: '',
    confirm_password: ''
  });

  const [profileErrors, setProfileErrors] = useState<Partial<ProfileFormData>>({});
  const [passwordErrors, setPasswordErrors] = useState<Partial<PasswordFormData>>({});

  useEffect(() => {
    loadUserData();
  }, [activeTab]);

  const loadUserData = async (): Promise<void> => {
    setIsLoading(true);
    try {
      if (activeTab === 'orders') {
        const response = await ordersAPI.getUserOrders();
        if (response.success && response.data) {
          setOrders(response.data.data || []);
        }
      } else if (activeTab === 'tickets') {
        const response = await ticketsAPI.getUserTickets();
        if (response.success && response.data) {
          setTickets(response.data.data || []);
        }
      }
    } catch (error) {
      console.error('Failed to load user data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleProfileSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    
    const errors: Partial<ProfileFormData> = {};
    if (!profileForm.first_name?.trim()) errors.first_name = 'First name is required';
    if (!profileForm.last_name?.trim()) errors.last_name = 'Last name is required';
    if (profileForm.phone && !/^\+?[\d\s-()]+$/.test(profileForm.phone)) {
      errors.phone = 'Please enter a valid phone number';
    }

    setProfileErrors(errors);
    if (Object.keys(errors).length > 0) return;

    setIsUpdating(true);
    try {
      const updateData: UpdateProfileData = {
        first_name: profileForm.first_name,
        last_name: profileForm.last_name,
        phone: profileForm.phone
      };

      const response = await userAPI.updateProfile(updateData);
      if (response.success && response.data) {
        updateUser(response.data);
        alert('Profile updated successfully!');
      } else {
        throw new Error(response.error || 'Failed to update profile');
      }
    } catch (error) {
      console.error('Failed to update profile:', error);
      alert('Failed to update profile. Please try again.');
    } finally {
      setIsUpdating(false);
    }
  };

  const handlePasswordSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    
    const errors: Partial<PasswordFormData> = {};
    if (!passwordForm.current_password) errors.current_password = 'Current password is required';
    if (!passwordForm.new_password) errors.new_password = 'New password is required';
    if (passwordForm.new_password.length < 6) errors.new_password = 'Password must be at least 6 characters';
    if (passwordForm.new_password !== passwordForm.confirm_password) {
      errors.confirm_password = 'Passwords do not match';
    }

    setPasswordErrors(errors);
    if (Object.keys(errors).length > 0) return;

    setIsUpdating(true);
    try {
      const response = await userAPI.changePassword(passwordForm);
      if (response.success) {
        setPasswordForm({
          current_password: '',
          new_password: '',
          confirm_password: ''
        });
        alert('Password changed successfully!');
      } else {
        throw new Error(response.error || 'Failed to change password');
      }
    } catch (error) {
      console.error('Failed to change password:', error);
      alert('Failed to change password. Please try again.');
    } finally {
      setIsUpdating(false);
    }
  };

  const getOrderStatusBadge = (status: string) => {
    const statusConfig = {
      'paid': { label: 'Paid', variant: 'default' as const },
      'confirmed': { label: 'Confirmed', variant: 'default' as const },
      'pending': { label: 'Pending', variant: 'secondary' as const },
      'cancelled': { label: 'Cancelled', variant: 'destructive' as const },
      'expired': { label: 'Expired', variant: 'destructive' as const }
    };
    
    const config = statusConfig[status as keyof typeof statusConfig] || { label: status, variant: 'outline' as const };
    return <Badge variant={config.variant}>{config.label}</Badge>;
  };

  const getTicketStatusIcon = (status: string) => {
    switch (status) {
      case 'active':
        return <CheckCircle className="h-5 w-5 text-green-500" />;
      case 'redeemed':
        return <CheckCircle className="h-5 w-5 text-blue-500" />;
      case 'voided':
        return <XCircle className="h-5 w-5 text-red-500" />;
      default:
        return <TicketIcon className="h-5 w-5 text-gray-500" />;
    }
  };

  const tabs = [
    { id: 'profile' as const, label: 'Profile', icon: User },
    { id: 'orders' as const, label: 'Orders', icon: Calendar },
    { id: 'tickets' as const, label: 'Tickets', icon: TicketIcon },
    { id: 'settings' as const, label: 'Settings', icon: Settings }
  ];

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8"
        >
          <h1 className="text-3xl font-bold text-gray-900 mb-2">My Account</h1>
          <p className="text-gray-600">Manage your profile, orders, and tickets</p>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          {/* Sidebar */}
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.1 }}
            className="lg:col-span-1"
          >
            <Card>
              <CardContent className="p-0">
                <nav className="space-y-1">
                  {tabs.map((tab) => (
                    <button
                      key={tab.id}
                      onClick={() => setActiveTab(tab.id)}
                      className={`w-full flex items-center px-4 py-3 text-left hover:bg-gray-50 transition-colors ${
                        activeTab === tab.id
                          ? 'bg-blue-50 text-blue-700 border-r-2 border-blue-700'
                          : 'text-gray-700'
                      }`}
                    >
                      <tab.icon className="h-5 w-5 mr-3" />
                      {tab.label}
                    </button>
                  ))}
                </nav>
              </CardContent>
            </Card>
          </motion.div>

          {/* Main Content */}
          <div className="lg:col-span-3">
            {/* Profile Tab */}
            {activeTab === 'profile' && (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
              >
                <Card>
                  <CardHeader>
                    <CardTitle>Profile Information</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <form onSubmit={handleProfileSubmit} className="space-y-4">
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                          <Label htmlFor="firstName">First Name *</Label>
                          <Input
                            id="firstName"
                            value={profileForm.first_name || ''}
                            onChange={(e) => setProfileForm(prev => ({ ...prev, first_name: e.target.value }))}
                            className={profileErrors.first_name ? 'border-red-500' : ''}
                          />
                          {profileErrors.first_name && (
                            <p className="text-sm text-red-600 mt-1">{profileErrors.first_name}</p>
                          )}
                        </div>
                        
                        <div>
                          <Label htmlFor="lastName">Last Name *</Label>
                          <Input
                            id="lastName"
                            value={profileForm.last_name || ''}
                            onChange={(e) => setProfileForm(prev => ({ ...prev, last_name: e.target.value }))}
                            className={profileErrors.last_name ? 'border-red-500' : ''}
                          />
                          {profileErrors.last_name && (
                            <p className="text-sm text-red-600 mt-1">{profileErrors.last_name}</p>
                          )}
                        </div>
                      </div>
                      
                      <div>
                        <Label htmlFor="email">Email</Label>
                        <Input
                          id="email"
                          type="email"
                          value={profileForm.email}
                          disabled
                          className="bg-gray-50"
                        />
                        <p className="text-sm text-gray-500 mt-1">Email cannot be changed</p>
                      </div>
                      
                      <div>
                        <Label htmlFor="phone">Phone Number</Label>
                        <Input
                          id="phone"
                          type="tel"
                          value={profileForm.phone || ''}
                          onChange={(e) => setProfileForm(prev => ({ ...prev, phone: e.target.value }))}
                          className={profileErrors.phone ? 'border-red-500' : ''}
                        />
                        {profileErrors.phone && (
                          <p className="text-sm text-red-600 mt-1">{profileErrors.phone}</p>
                        )}
                      </div>
                      
                      <Button type="submit" disabled={isUpdating}>
                        {isUpdating ? (
                          <>
                            <LoadingSpinner />
                            Updating...
                          </>
                        ) : (
                          'Update Profile'
                        )}
                      </Button>
                    </form>
                  </CardContent>
                </Card>
              </motion.div>
            )}

            {/* Orders Tab */}
            {activeTab === 'orders' && (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
              >
                <Card>
                  <CardHeader>
                    <CardTitle>Order History</CardTitle>
                  </CardHeader>
                  <CardContent>
                    {isLoading ? (
                      <div className="flex justify-center py-8">
                        <LoadingSpinner />
                      </div>
                    ) : orders.length > 0 ? (
                      <div className="space-y-4">
                        {orders.map((order) => (
                          <div key={order.id} className="border rounded-lg p-4">
                            <div className="flex justify-between items-start mb-3">
                              <div>
                                <h3 className="font-semibold">Order #{order.code}</h3>
                                <p className="text-sm text-gray-600">
                                  {formatDateTime(order.created_at)}
                                </p>
                                {order.event && (
                                  <p className="text-sm text-gray-600">{order.event.name}</p>
                                )}
                              </div>
                              <div className="text-right">
                                {getOrderStatusBadge(order.status)}
                                <p className="font-semibold mt-1">
                                  {formatCurrency(order.total_amount)}
                                </p>
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>
                    ) : (
                      <div className="text-center py-8">
                        <Calendar className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                        <p className="text-gray-600">No orders found</p>
                      </div>
                    )}
                  </CardContent>
                </Card>
              </motion.div>
            )}

            {/* Tickets Tab */}
            {activeTab === 'tickets' && (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
              >
                <Card>
                  <CardHeader>
                    <CardTitle>My Tickets</CardTitle>
                  </CardHeader>
                  <CardContent>
                    {isLoading ? (
                      <div className="flex justify-center py-8">
                        <LoadingSpinner />
                      </div>
                    ) : tickets.length > 0 ? (
                      <div className="space-y-4">
                        {tickets.map((ticket) => (
                          <div key={ticket.id} className="border rounded-lg p-4">
                            <div className="flex justify-between items-start">
                              <div>
                                <h3 className="font-semibold">Ticket #{ticket.serial_number}</h3>
                                <p className="text-sm text-gray-600">
                                  Created: {formatDateTime(ticket.created_at)}
                                </p>
                                {ticket.redeemed_at && (
                                  <p className="text-sm text-gray-600">
                                    Redeemed: {formatDateTime(ticket.redeemed_at)}
                                  </p>
                                )}
                              </div>
                              <div className="flex items-center">
                                {getTicketStatusIcon(ticket.status)}
                                <span className="ml-2 text-sm capitalize">{ticket.status}</span>
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>
                    ) : (
                      <div className="text-center py-8">
                        <TicketIcon className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                        <p className="text-gray-600">No tickets found</p>
                      </div>
                    )}
                  </CardContent>
                </Card>
              </motion.div>
            )}

            {/* Settings Tab */}
            {activeTab === 'settings' && (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
              >
                <Card>
                  <CardHeader>
                    <CardTitle>Change Password</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <form onSubmit={handlePasswordSubmit} className="space-y-4">
                      <div>
                        <Label htmlFor="currentPassword">Current Password *</Label>
                        <div className="relative">
                          <Input
                            id="currentPassword"
                            type={showPassword ? 'text' : 'password'}
                            value={passwordForm.current_password}
                            onChange={(e) => setPasswordForm(prev => ({ ...prev, current_password: e.target.value }))}
                            className={passwordErrors.current_password ? 'border-red-500' : ''}
                          />
                          <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword)}
                            className="absolute right-3 top-1/2 transform -translate-y-1/2"
                          >
                            {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                          </button>
                        </div>
                        {passwordErrors.current_password && (
                          <p className="text-sm text-red-600 mt-1">{passwordErrors.current_password}</p>
                        )}
                      </div>
                      
                      <div>
                        <Label htmlFor="newPassword">New Password *</Label>
                        <Input
                          id="newPassword"
                          type="password"
                          value={passwordForm.new_password}
                          onChange={(e) => setPasswordForm(prev => ({ ...prev, new_password: e.target.value }))}
                          className={passwordErrors.new_password ? 'border-red-500' : ''}
                        />
                        {passwordErrors.new_password && (
                          <p className="text-sm text-red-600 mt-1">{passwordErrors.new_password}</p>
                        )}
                      </div>
                      
                      <div>
                        <Label htmlFor="confirmPassword">Confirm New Password *</Label>
                        <Input
                          id="confirmPassword"
                          type="password"
                          value={passwordForm.confirm_password}
                          onChange={(e) => setPasswordForm(prev => ({ ...prev, confirm_password: e.target.value }))}
                          className={passwordErrors.confirm_password ? 'border-red-500' : ''}
                        />
                        {passwordErrors.confirm_password && (
                          <p className="text-sm text-red-600 mt-1">{passwordErrors.confirm_password}</p>
                        )}
                      </div>
                      
                      <Button type="submit" disabled={isUpdating}>
                        {isUpdating ? (
                          <>
                            <LoadingSpinner />
                            Changing...
                          </>
                        ) : (
                          'Change Password'
                        )}
                      </Button>
                    </form>
                  </CardContent>
                </Card>
              </motion.div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProfilePage;

