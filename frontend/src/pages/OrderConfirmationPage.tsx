import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { ordersAPI, ticketsAPI } from '../services/api';
import { Order, Ticket } from '../types/api';
import { formatCurrency, formatDateTime } from '../lib/utils';
import { 
  CheckCircle, 
  XCircle, 
  Clock, 
  Download, 
  Mail, 
  Ticket as TicketIcon,
  Calendar,
  MapPin,
  User
} from 'lucide-react';
import { motion } from 'framer-motion';

const OrderConfirmationPage: React.FC = () => {
  const { orderId } = useParams<{ orderId: string }>();
  const [order, setOrder] = useState<Order | null>(null);
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isDownloading, setIsDownloading] = useState<boolean>(false);

  useEffect(() => {
    if (orderId) {
      loadOrderDetails();
    }
  }, [orderId]);

  const loadOrderDetails = async (): Promise<void> => {
    if (!orderId) return;
    
    try {
      const [orderResponse, ticketsResponse] = await Promise.all([
        ordersAPI.getOrder(orderId),
        ticketsAPI.getUserTickets({ order_id: orderId })
      ]);

      if (orderResponse.success && orderResponse.data) {
        setOrder(orderResponse.data);
      }

      if (ticketsResponse.success && ticketsResponse.data) {
        setTickets(ticketsResponse.data.data || []);
      }
    } catch (error) {
      console.error('Failed to load order details:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDownloadTickets = async (): Promise<void> => {
    if (!tickets.length) return;
    
    setIsDownloading(true);
    try {
      for (const ticket of tickets) {
        const response = await ticketsAPI.downloadTicket(ticket.id);
        if (response.success) {
          // Handle download - this would typically trigger a file download
          console.log('Ticket downloaded:', ticket.id);
        }
      }
    } catch (error) {
      console.error('Failed to download tickets:', error);
    } finally {
      setIsDownloading(false);
    }
  };

  const handleResendTickets = async (): Promise<void> => {
    if (!tickets.length) return;
    
    try {
      for (const ticket of tickets) {
        await ticketsAPI.resendTicket(ticket.id);
      }
      alert('Tickets have been resent to your email!');
    } catch (error) {
      console.error('Failed to resend tickets:', error);
      alert('Failed to resend tickets. Please try again.');
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'paid':
      case 'confirmed':
        return <CheckCircle className="h-8 w-8 text-green-500" />;
      case 'pending':
        return <Clock className="h-8 w-8 text-yellow-500" />;
      case 'cancelled':
      case 'expired':
        return <XCircle className="h-8 w-8 text-red-500" />;
      default:
        return <Clock className="h-8 w-8 text-gray-500" />;
    }
  };

  const getStatusMessage = (status: string) => {
    switch (status) {
      case 'paid':
      case 'confirmed':
        return {
          title: 'Payment Successful!',
          message: 'Your order has been confirmed and your tickets are ready.',
          color: 'text-green-600'
        };
      case 'pending':
        return {
          title: 'Payment Pending',
          message: 'We are processing your payment. Please wait a moment.',
          color: 'text-yellow-600'
        };
      case 'cancelled':
        return {
          title: 'Order Cancelled',
          message: 'This order has been cancelled.',
          color: 'text-red-600'
        };
      case 'expired':
        return {
          title: 'Order Expired',
          message: 'This order has expired. Please create a new order.',
          color: 'text-red-600'
        };
      default:
        return {
          title: 'Order Status Unknown',
          message: 'Please contact support for assistance.',
          color: 'text-gray-600'
        };
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  if (!order) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <XCircle className="h-16 w-16 text-red-500 mx-auto mb-4" />
          <h2 className="text-2xl font-semibold text-gray-600 mb-2">Order not found</h2>
          <p className="text-gray-500 mb-6">The order you're looking for doesn't exist or has been removed.</p>
          <Button asChild>
            <Link to="/events">Browse Events</Link>
          </Button>
        </div>
      </div>
    );
  }

  const statusInfo = getStatusMessage(order.status);

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Status Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-center mb-8"
        >
          <div className="flex justify-center mb-4">
            {getStatusIcon(order.status)}
          </div>
          <h1 className={`text-3xl font-bold mb-2 ${statusInfo.color}`}>
            {statusInfo.title}
          </h1>
          <p className="text-gray-600 text-lg">
            {statusInfo.message}
          </p>
          <p className="text-sm text-gray-500 mt-2">
            Order #{order.code}
          </p>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Order Details */}
          <div className="lg:col-span-2 space-y-6">
            {/* Event Information */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 }}
            >
              <Card>
                <CardHeader>
                  <CardTitle>Event Details</CardTitle>
                </CardHeader>
                <CardContent>
                  {order.event && (
                    <div className="space-y-3">
                      <h3 className="text-xl font-semibold">{order.event.name}</h3>
                      <div className="space-y-2 text-gray-600">
                        <div className="flex items-center">
                          <Calendar className="h-5 w-5 mr-3" />
                          <span>{formatDateTime(order.event.event_date)}</span>
                        </div>
                        <div className="flex items-center">
                          <MapPin className="h-5 w-5 mr-3" />
                          <span>{order.event.venue_name}, {order.event.venue_city}</span>
                        </div>
                      </div>
                    </div>
                  )}
                </CardContent>
              </Card>
            </motion.div>

            {/* Order Items */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.2 }}
            >
              <Card>
                <CardHeader>
                  <CardTitle>Order Items</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {order.order_lines?.map((line) => (
                      <div key={line.id} className="flex justify-between items-center p-4 border rounded-lg">
                        <div>
                          <h4 className="font-semibold">
                            {line.ticket_tier?.name || 'Ticket'}
                          </h4>
                          <p className="text-sm text-gray-600">
                            Quantity: {line.quantity}
                          </p>
                          <p className="text-sm text-gray-600">
                            {formatCurrency(line.unit_price)} each
                          </p>
                        </div>
                        <div className="text-right">
                          <p className="font-semibold">
                            {formatCurrency(line.total_price)}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </motion.div>

            {/* Customer Information */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
            >
              <Card>
                <CardHeader>
                  <CardTitle>Customer Information</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    <div className="flex items-center">
                      <User className="h-5 w-5 mr-3" />
                      <span>{order.customer_first_name} {order.customer_last_name}</span>
                    </div>
                    <div className="flex items-center">
                      <Mail className="h-5 w-5 mr-3" />
                      <span>{order.customer_email}</span>
                    </div>
                    {order.customer_phone && (
                      <div className="flex items-center">
                        <span className="w-5 h-5 mr-3 text-center">📱</span>
                        <span>{order.customer_phone}</span>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            </motion.div>
          </div>

          {/* Order Summary & Actions */}
          <div className="space-y-6">
            {/* Order Summary */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.4 }}
            >
              <Card>
                <CardHeader>
                  <CardTitle>Order Summary</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span>Subtotal</span>
                      <span>{formatCurrency(order.total_amount)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Service Fee</span>
                      <span>{formatCurrency(0)}</span>
                    </div>
                    <hr />
                    <div className="flex justify-between text-lg font-bold">
                      <span>Total</span>
                      <span>{formatCurrency(order.total_amount)}</span>
                    </div>
                  </div>
                  
                  <div className="pt-4">
                    <Badge variant={order.status === 'paid' ? 'default' : 'secondary'}>
                      {order.status.toUpperCase()}
                    </Badge>
                  </div>
                </CardContent>
              </Card>
            </motion.div>

            {/* Actions */}
            {(order.status === 'paid' || order.status === 'confirmed') && (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.5 }}
              >
                <Card>
                  <CardHeader>
                    <CardTitle>Your Tickets</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="text-center">
                      <TicketIcon className="h-12 w-12 text-green-500 mx-auto mb-2" />
                      <p className="text-sm text-gray-600 mb-4">
                        {tickets.length} ticket{tickets.length > 1 ? 's' : ''} ready
                      </p>
                    </div>
                    
                    <div className="space-y-2">
                      <Button
                        onClick={handleDownloadTickets}
                        disabled={isDownloading}
                        className="w-full"
                      >
                        {isDownloading ? (
                          <>
                            <LoadingSpinner />
                            Downloading...
                          </>
                        ) : (
                          <>
                            <Download className="h-4 w-4 mr-2" />
                            Download Tickets
                          </>
                        )}
                      </Button>
                      
                      <Button
                        variant="outline"
                        onClick={handleResendTickets}
                        className="w-full"
                      >
                        <Mail className="h-4 w-4 mr-2" />
                        Resend to Email
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              </motion.div>
            )}

            {/* Navigation */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.6 }}
            >
              <div className="space-y-2">
                <Button asChild variant="outline" className="w-full">
                  <Link to="/events">Browse More Events</Link>
                </Button>
                <Button asChild variant="ghost" className="w-full">
                  <Link to="/profile">View Order History</Link>
                </Button>
              </div>
            </motion.div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default OrderConfirmationPage;

