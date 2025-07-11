import { Heart } from "lucide-react";

export const Footer = () => {
  return (
    <footer className="bg-card border-t border-border mt-auto">
      <div className="container mx-auto px-4 py-6">
        <div className="text-center text-muted-foreground">
          <p className="flex items-center justify-center gap-1 text-sm">
            Made with <Heart className="h-4 w-4 text-red-500 fill-current" /> for English learning
          </p>
          <p className="text-xs mt-2">
            © 2024 Essay Recite. Practice makes perfect.
          </p>
        </div>
      </div>
    </footer>
  );
};