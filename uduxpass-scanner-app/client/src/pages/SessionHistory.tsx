/**
 * Design Philosophy: Professional Event Tech - Clean list view with session cards
 * Layout: List of session cards with statistics and status badges
 */

import { useState, useEffect } from 'react';
import { useLocation } from 'wouter';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { scannerApi, ScanningSession } from '@/lib/api';
import { ArrowLeft, Calendar, MapPin, Scan, CheckCircle, XCircle, Loader2 } from 'lucide-react';
import { toast } from 'sonner';

export default function SessionHistory() {
  const [, setLocation] = useLocation();
  const [sessions, setSessions] = useState<ScanningSession[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadSessions();
  }, []);

  const loadSessions = async () => {
    try {
      const sessionsData = await scannerApi.getAllSessions();
      setSessions(sessionsData);
    } catch (error) {
      console.error('Failed to load sessions:', error);
      toast.error('Failed to load session history');
    } finally {
      setIsLoading(false);
    }
  };

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getStatusBadge = (status: string) => {
    if (status === 'active') {
      return (
        <Badge className="bg-blue-600 text-white hover:bg-blue-700">
          Active
        </Badge>
      );
    }
    return (
      <Badge variant="secondary" className="bg-green-100 text-green-700">
        Completed
      </Badge>
    );
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Top Bar */}
      <header className="bg-primary text-primary-foreground px-4 py-3 flex items-center gap-3 sticky top-0 z-10">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setLocation('/dashboard')}
          className="text-primary-foreground hover:bg-primary-foreground/10"
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <h1 className="font-semibold text-lg">Session History</h1>
      </header>

      {/* Content */}
      <main className="p-4">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        ) : sessions.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <Scan className="h-12 w-12 text-muted-foreground mx-auto mb-3" />
              <p className="text-muted-foreground mb-1">No sessions found</p>
              <p className="text-sm text-muted-foreground">
                Start a new session to begin scanning tickets
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-3">
            {sessions.map((session) => (
              <Card key={session.id} className="hover:shadow-md transition-shadow">
                <CardContent className="p-4">
                  {/* Header */}
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex items-start gap-2 flex-1">
                      <Calendar className="h-5 w-5 text-primary mt-0.5 flex-shrink-0" />
                      <div>
                        <h3 className="font-semibold text-foreground">
                          {session.event?.name || 'Event'}
                        </h3>
                        <p className="text-sm text-muted-foreground">
                          {formatDateTime(session.start_time)}
                          {session.end_time && ` - ${formatDateTime(session.end_time)}`}
                        </p>
                      </div>
                    </div>
                    {getStatusBadge(session.status)}
                  </div>

                  {/* Location */}
                  <div className="flex items-center gap-2 mb-3">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <p className="text-sm text-muted-foreground">
                      {session.location}
                    </p>
                  </div>

                  {/* Statistics */}
                  <div className="flex items-center gap-4 text-sm">
                    <div className="flex items-center gap-1.5">
                      <Scan className="h-4 w-4 text-muted-foreground" />
                      <span className="font-medium">0</span>
                      <span className="text-muted-foreground">scanned</span>
                    </div>
                    <div className="flex items-center gap-1.5">
                      <CheckCircle className="h-4 w-4 text-green-600" />
                      <span className="font-medium text-green-600">0</span>
                      <span className="text-muted-foreground">valid</span>
                    </div>
                    <div className="flex items-center gap-1.5">
                      <XCircle className="h-4 w-4 text-red-600" />
                      <span className="font-medium text-red-600">0</span>
                      <span className="text-muted-foreground">invalid</span>
                    </div>
                  </div>

                  {/* Notes */}
                  {session.notes && (
                    <div className="mt-3 pt-3 border-t border-border">
                      <p className="text-sm text-muted-foreground italic">
                        {session.notes}
                      </p>
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
