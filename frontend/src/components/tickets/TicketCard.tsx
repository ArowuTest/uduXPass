import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { TicketQRCode } from './TicketQRCode';
import { Calendar, MapPin, Tag, Hash } from 'lucide-react';
import { format } from 'date-fns';

interface TicketCardProps {
  ticket: {
    id: string;
    serial_number: string;
    qr_code_data: string;
    status: 'active' | 'redeemed' | 'cancelled';
    redeemed_at?: string;
    created_at: string;
    event?: {
      title: string;
      venue: string;
      start_date: string;
    };
    tier?: {
      name: string;
      price: number;
    };
  };
  expanded?: boolean;
}

export function TicketCard({ ticket, expanded = false }: TicketCardProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800 border-green-200';
      case 'redeemed':
        return 'bg-gray-100 text-gray-800 border-gray-200';
      case 'cancelled':
        return 'bg-red-100 text-red-800 border-red-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'active':
        return 'Valid';
      case 'redeemed':
        return 'Used';
      case 'cancelled':
        return 'Cancelled';
      default:
        return status;
    }
  };

  return (
    <Card className="overflow-hidden">
      <CardHeader className="bg-gradient-to-r from-blue-50 to-indigo-50 border-b">
        <div className="flex items-start justify-between">
          <div className="space-y-1">
            <CardTitle className="text-xl">
              {ticket.event?.title || 'Event Ticket'}
            </CardTitle>
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Hash className="h-4 w-4" />
              <span className="font-mono">{ticket.serial_number}</span>
            </div>
          </div>
          <Badge className={getStatusColor(ticket.status)}>
            {getStatusText(ticket.status)}
          </Badge>
        </div>
      </CardHeader>

      <CardContent className="p-6">
        <div className="grid md:grid-cols-2 gap-6">
          {/* Ticket Details */}
          <div className="space-y-4">
            <div>
              <h3 className="font-semibold mb-3 text-lg">Ticket Details</h3>
              <div className="space-y-3">
                {ticket.tier && (
                  <div className="flex items-center gap-2 text-sm">
                    <Tag className="h-4 w-4 text-muted-foreground" />
                    <span className="text-muted-foreground">Tier:</span>
                    <span className="font-medium">{ticket.tier.name}</span>
                  </div>
                )}

                {ticket.event?.venue && (
                  <div className="flex items-center gap-2 text-sm">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <span className="text-muted-foreground">Venue:</span>
                    <span className="font-medium">{ticket.event.venue}</span>
                  </div>
                )}

                {ticket.event?.start_date && (
                  <div className="flex items-center gap-2 text-sm">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <span className="text-muted-foreground">Date:</span>
                    <span className="font-medium">
                      {format(new Date(ticket.event.start_date), 'PPP')}
                    </span>
                  </div>
                )}

                {ticket.redeemed_at && (
                  <div className="flex items-center gap-2 text-sm">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <span className="text-muted-foreground">Scanned:</span>
                    <span className="font-medium">
                      {format(new Date(ticket.redeemed_at), 'PPP p')}
                    </span>
                  </div>
                )}
              </div>
            </div>

            {ticket.tier?.price && (
              <div className="pt-4 border-t">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Price Paid</span>
                  <span className="text-lg font-bold">
                    ₦{ticket.tier.price.toLocaleString()}
                  </span>
                </div>
              </div>
            )}
          </div>

          {/* QR Code */}
          <div className="flex flex-col items-center justify-center">
            {ticket.status === 'active' ? (
              <>
                <TicketQRCode
                  qrCodeData={ticket.qr_code_data}
                  ticketSerial={ticket.serial_number}
                  size={expanded ? 300 : 200}
                  showActions={expanded}
                />
                <p className="text-xs text-muted-foreground text-center mt-2 max-w-xs">
                  Present this QR code at the event entrance for scanning
                </p>
              </>
            ) : (
              <div className="text-center p-8 bg-gray-50 rounded-lg">
                <p className="text-sm text-muted-foreground">
                  {ticket.status === 'redeemed' 
                    ? 'This ticket has been used'
                    : 'This ticket is no longer valid'}
                </p>
              </div>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
