/**
 * Design Philosophy: Professional Event Tech - Clear status overview, quick actions
 * Layout: Top bar with scanner info, active session card, statistics grid, action buttons
 * Colors: Blue primary, green success, red error
 */

import { useState, useEffect } from 'react';
import { useLocation } from 'wouter';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useAuth } from '@/contexts/AuthContext';
import { scannerApi, ScanningSession, SessionStats } from '@/lib/api';
import { LogOut, QrCode, PlusCircle, Calendar, MapPin, CheckCircle, XCircle, Scan, History } from 'lucide-react';
import { toast } from 'sonner';

export default function Dashboard() {
  const [, setLocation] = useLocation();
  const { scanner, logout } = useAuth();
  const [activeSession, setActiveSession] = useState<ScanningSession | null>(null);
  const [stats, setStats] = useState<SessionStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadActiveSession();
  }, []);

  const loadActiveSession = async () => {
    try {
      const sessions = await scannerApi.getActiveSessions();
      if (sessions.length > 0) {
        const session = sessions[0];
        setActiveSession(session);
        
        // Load stats for active session
        const sessionStats = await scannerApi.getSessionStats(session.id);
        setStats(sessionStats);
      }
    } catch (error) {
      console.error('Failed to load active session:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    toast.success('Logged out successfully');
    setLocation('/login');
  };

  const handleScanTicket = () => {
    if (!activeSession) {
      toast.error('Please start a scanning session first');
      return;
    }
    setLocation('/scan');
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Top Bar */}
      <header className="bg-card border-b border-border px-4 py-3 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-primary rounded-lg flex items-center justify-center">
            <span className="text-primary-foreground font-bold text-lg">U</span>
          </div>
          <div>
            <h1 className="font-semibold text-foreground">uduXPass</h1>
            <p className="text-xs text-muted-foreground">{scanner?.name}</p>
          </div>
        </div>
        <Button
          variant="ghost"
          size="sm"
          onClick={handleLogout}
          className="text-muted-foreground hover:text-foreground"
        >
          <LogOut className="h-4 w-4 mr-2" />
          Logout
        </Button>
      </header>

      {/* Main Content */}
      <main className="p-4 space-y-4">
        {/* Active Session Card */}
        {activeSession && (
          <Card className="border-primary/20 bg-primary/5">
            <CardHeader className="pb-3">
              <CardTitle className="text-lg flex items-center gap-2 text-primary">
                <Scan className="h-5 w-5" />
                Active Scanning Session
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <div className="flex items-start gap-2">
                <Calendar className="h-4 w-4 text-muted-foreground mt-0.5" />
                <div>
                  <p className="font-semibold text-foreground">
                    {activeSession.event?.name || 'Event'}
                  </p>
                  <p className="text-sm text-muted-foreground">
                    {activeSession.event?.start_time && new Date(activeSession.event.start_time).toLocaleDateString()}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <MapPin className="h-4 w-4 text-muted-foreground" />
                <p className="text-sm text-muted-foreground">
                  {activeSession.location || 'Location not specified'}
                </p>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Statistics Grid */}
        {stats && (
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-2xl font-bold text-foreground">{stats.total_scanned}</p>
                    <p className="text-sm text-muted-foreground">Tickets Scanned</p>
                  </div>
                  <Scan className="h-8 w-8 text-primary/50" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-2xl font-bold text-green-600">{stats.valid_tickets}</p>
                    <p className="text-sm text-muted-foreground">Valid Tickets</p>
                  </div>
                  <CheckCircle className="h-8 w-8 text-green-600/50" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-2xl font-bold text-red-600">{stats.invalid_tickets}</p>
                    <p className="text-sm text-muted-foreground">Invalid Tickets</p>
                  </div>
                  <XCircle className="h-8 w-8 text-red-600/50" />
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {/* No Active Session Message */}
        {!activeSession && !isLoading && (
          <Card>
            <CardContent className="py-8 text-center">
              <Scan className="h-12 w-12 text-muted-foreground mx-auto mb-3" />
              <p className="text-muted-foreground mb-1">No active scanning session</p>
              <p className="text-sm text-muted-foreground">Start a new session to begin scanning tickets</p>
            </CardContent>
          </Card>
        )}

        {/* Action Buttons */}
        <div className="space-y-3 pt-2">
          <Button
            className="w-full h-14 text-base font-semibold"
            onClick={() => setLocation('/create-session')}
          >
            <PlusCircle className="mr-2 h-5 w-5" />
            Start New Session
          </Button>

          <Button
            className="w-full h-14 text-base font-semibold bg-green-600 hover:bg-green-700 text-white"
            onClick={handleScanTicket}
            disabled={!activeSession}
          >
            <QrCode className="mr-2 h-5 w-5" />
            Scan Ticket
          </Button>

          <Button
            variant="outline"
            className="w-full h-12 text-base"
            onClick={() => setLocation('/history')}
          >
            <History className="mr-2 h-5 w-5" />
            Session History
          </Button>
        </div>
      </main>
    </div>
  );
}
