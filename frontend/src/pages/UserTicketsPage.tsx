import { useState, useEffect } from 'react';
import { TicketCard } from '@/components/tickets/TicketCard';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Search, Ticket, AlertCircle } from 'lucide-react';
import { toast } from 'sonner';
import api from '@/services/api';

interface Ticket {
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
}

export default function UserTicketsPage() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedTab, setSelectedTab] = useState('all');

  useEffect(() => {
    loadTickets();
  }, []);

  const loadTickets = async () => {
    try {
      setLoading(true);
      const response = await api.get('/user/tickets');
      setTickets(response.data.data.tickets || []);
    } catch (error: any) {
      console.error('Failed to load tickets:', error);
      toast.error(error.response?.data?.message || 'Failed to load tickets');
    } finally {
      setLoading(false);
    }
  };

  const filterTickets = (status?: string) => {
    let filtered = tickets;

    // Filter by status
    if (status && status !== 'all') {
      filtered = filtered.filter(t => t.status === status);
    }

    // Filter by search term
    if (searchTerm) {
      filtered = filtered.filter(t =>
        t.serial_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
        t.event?.title.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    return filtered;
  };

  const getTicketCount = (status?: string) => {
    if (!status || status === 'all') return tickets.length;
    return tickets.filter(t => t.status === status).length;
  };

  if (loading) {
    return (
      <div className="container py-8">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading your tickets...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="container py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">My Tickets</h1>
        <p className="text-muted-foreground">
          View and manage your event tickets
        </p>
      </div>

      {/* Search */}
      <div className="mb-6">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search by ticket number or event name..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-10"
          />
        </div>
      </div>

      {/* Tabs */}
      <Tabs value={selectedTab} onValueChange={setSelectedTab} className="mb-6">
        <TabsList>
          <TabsTrigger value="all">
            All ({getTicketCount('all')})
          </TabsTrigger>
          <TabsTrigger value="active">
            Active ({getTicketCount('active')})
          </TabsTrigger>
          <TabsTrigger value="redeemed">
            Used ({getTicketCount('redeemed')})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="all" className="mt-6">
          <TicketList tickets={filterTickets('all')} />
        </TabsContent>

        <TabsContent value="active" className="mt-6">
          <TicketList tickets={filterTickets('active')} />
        </TabsContent>

        <TabsContent value="redeemed" className="mt-6">
          <TicketList tickets={filterTickets('redeemed')} />
        </TabsContent>
      </Tabs>
    </div>
  );
}

function TicketList({ tickets }: { tickets: Ticket[] }) {
  if (tickets.length === 0) {
    return (
      <div className="text-center py-12 bg-muted/30 rounded-lg">
        <Ticket className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
        <h3 className="text-lg font-semibold mb-2">No tickets found</h3>
        <p className="text-muted-foreground mb-4">
          You don't have any tickets matching this filter
        </p>
        <Button variant="outline" onClick={() => window.location.href = '/events'}>
          Browse Events
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {tickets.map((ticket) => (
        <TicketCard key={ticket.id} ticket={ticket} expanded={false} />
      ))}
    </div>
  );
}
