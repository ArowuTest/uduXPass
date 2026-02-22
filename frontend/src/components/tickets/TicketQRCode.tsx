import { QRCodeSVG } from 'qrcode.react';
import { Button } from '@/components/ui/button';
import { Download, Share2 } from 'lucide-react';
import { toast } from 'sonner';

interface TicketQRCodeProps {
  qrCodeData: string;
  ticketSerial: string;
  size?: number;
  showActions?: boolean;
}

export function TicketQRCode({ 
  qrCodeData, 
  ticketSerial, 
  size = 256,
  showActions = true 
}: TicketQRCodeProps) {
  
  const downloadQRCode = () => {
    try {
      // Get the SVG element
      const svg = document.getElementById(`qr-${ticketSerial}`);
      if (!svg) {
        toast.error('Failed to find QR code');
        return;
      }

      // Convert SVG to canvas
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      if (!ctx) {
        toast.error('Failed to create canvas');
        return;
      }

      const svgData = new XMLSerializer().serializeToString(svg);
      const img = new Image();
      
      img.onload = () => {
        canvas.width = size;
        canvas.height = size;
        ctx.drawImage(img, 0, 0);
        
        // Download as PNG
        canvas.toBlob((blob) => {
          if (!blob) {
            toast.error('Failed to generate image');
            return;
          }
          
          const url = URL.createObjectURL(blob);
          const link = document.createElement('a');
          link.href = url;
          link.download = `ticket-${ticketSerial}.png`;
          link.click();
          URL.revokeObjectURL(url);
          
          toast.success('QR code downloaded');
        });
      };
      
      img.src = 'data:image/svg+xml;base64,' + btoa(svgData);
    } catch (error) {
      console.error('Failed to download QR code:', error);
      toast.error('Failed to download QR code');
    }
  };

  const shareQRCode = async () => {
    try {
      if (navigator.share) {
        await navigator.share({
          title: `Ticket ${ticketSerial}`,
          text: `My ticket QR code: ${ticketSerial}`,
          url: window.location.href,
        });
        toast.success('Shared successfully');
      } else {
        // Fallback: copy to clipboard
        await navigator.clipboard.writeText(qrCodeData);
        toast.success('QR code data copied to clipboard');
      }
    } catch (error) {
      console.error('Failed to share:', error);
      toast.error('Failed to share');
    }
  };

  return (
    <div className="flex flex-col items-center gap-4">
      <div className="bg-white p-4 rounded-lg shadow-sm border">
        <QRCodeSVG
          id={`qr-${ticketSerial}`}
          value={qrCodeData}
          size={size}
          level="H" // High error correction
          includeMargin={true}
        />
      </div>
      
      {showActions && (
        <div className="flex gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={downloadQRCode}
            className="gap-2"
          >
            <Download className="h-4 w-4" />
            Download
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={shareQRCode}
            className="gap-2"
          >
            <Share2 className="h-4 w-4" />
            Share
          </Button>
        </div>
      )}
    </div>
  );
}
