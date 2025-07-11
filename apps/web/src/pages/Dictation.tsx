import { useState, useEffect, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { Progress } from "@/components/ui/progress";
import { ArrowLeft, Play, Pause, RotateCcw, Check, X } from "lucide-react";
import { useToast } from "@/hooks/use-toast";

interface Essay {
  id: string;
  title: string;
  content: string;
  createdAt: string;
  audioUrl?: string;
  imageUrl?: string;
  wordCount: number;
}

const Dictation = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const [essay, setEssay] = useState<Essay | null>(null);
  const [sentences, setSentences] = useState<string[]>([]);
  const [currentSentence, setCurrentSentence] = useState(0);
  const [userInput, setUserInput] = useState("");
  const [isPlaying, setIsPlaying] = useState(false);
  const [isRevealed, setIsRevealed] = useState(false);
  const [score, setScore] = useState({ correct: 0, total: 0 });
  const audioRef = useRef<HTMLAudioElement>(null);

  useEffect(() => {
    if (id) {
      const essays = JSON.parse(localStorage.getItem('essays') || '[]');
      const foundEssay = essays.find((e: Essay) => e.id === id);
      if (foundEssay) {
        setEssay(foundEssay);
        const sentenceArray = foundEssay.content
          .split(/[.!?]+/)
          .filter(s => s.trim().length > 0)
          .map(s => s.trim() + '.');
        setSentences(sentenceArray);
      }
    }
  }, [id]);

  const handlePlayPause = () => {
    setIsPlaying(!isPlaying);
    // In a real app, this would control actual audio playback
    toast({
      title: isPlaying ? "Paused" : "Playing",
      description: `Sentence ${currentSentence + 1}`,
    });
  };

  const checkAnswer = () => {
    const correctSentence = sentences[currentSentence];
    const similarity = calculateSimilarity(userInput.trim(), correctSentence);
    const isCorrect = similarity > 0.8;
    
    setScore(prev => ({
      correct: prev.correct + (isCorrect ? 1 : 0),
      total: prev.total + 1
    }));

    setIsRevealed(true);
    
    toast({
      title: isCorrect ? "Correct!" : "Try Again",
      description: isCorrect 
        ? "Great job!" 
        : `Similarity: ${Math.round(similarity * 100)}%`,
    });
  };

  const calculateSimilarity = (input: string, target: string): number => {
    const inputWords = input.toLowerCase().split(/\s+/);
    const targetWords = target.toLowerCase().split(/\s+/);
    const matchingWords = inputWords.filter(word => targetWords.includes(word));
    return matchingWords.length / Math.max(inputWords.length, targetWords.length);
  };

  const nextSentence = () => {
    if (currentSentence < sentences.length - 1) {
      setCurrentSentence(prev => prev + 1);
      setUserInput("");
      setIsRevealed(false);
      setIsPlaying(false);
    } else {
      toast({
        title: "Dictation Complete!",
        description: `Score: ${score.correct + (isRevealed ? 1 : 0)}/${score.total + 1}`,
      });
    }
  };

  const resetDictation = () => {
    setCurrentSentence(0);
    setUserInput("");
    setIsRevealed(false);
    setIsPlaying(false);
    setScore({ correct: 0, total: 0 });
  };

  if (!essay) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p className="text-center text-muted-foreground">Essay not found</p>
      </div>
    );
  }

  const progress = ((currentSentence + 1) / sentences.length) * 100;

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-6">
        <Button 
          variant="outline" 
          onClick={() => navigate(`/essay/${id}`)}
          className="mb-4"
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Essay
        </Button>
        <h1 className="text-3xl font-bold text-foreground mb-2">Dictation Training</h1>
        <p className="text-muted-foreground">{essay.title}</p>
      </div>

      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Progress</CardTitle>
              <div className="text-sm text-muted-foreground">
                Sentence {currentSentence + 1} of {sentences.length}
              </div>
            </div>
            <Progress value={progress} className="mt-2" />
          </CardHeader>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Audio Controls</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-4">
              <Button onClick={handlePlayPause} variant="default">
                {isPlaying ? <Pause className="h-4 w-4" /> : <Play className="h-4 w-4" />}
                {isPlaying ? "Pause" : "Play"}
              </Button>
              <Button onClick={resetDictation} variant="outline">
                <RotateCcw className="h-4 w-4 mr-2" />
                Reset
              </Button>
              <div className="text-sm text-muted-foreground">
                Score: {score.correct}/{score.total}
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Type What You Hear</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <Textarea
              placeholder="Type the sentence you heard..."
              value={userInput}
              onChange={(e) => setUserInput(e.target.value)}
              disabled={isRevealed}
              className="min-h-[100px]"
            />
            
            {isRevealed && (
              <div className="space-y-3">
                <div className="p-3 bg-muted rounded-lg">
                  <p className="text-sm font-medium mb-1">Correct sentence:</p>
                  <p className="text-foreground">{sentences[currentSentence]}</p>
                </div>
                <div className="p-3 bg-accent/10 rounded-lg">
                  <p className="text-sm font-medium mb-1">Your input:</p>
                  <p className="text-foreground">{userInput}</p>
                </div>
              </div>
            )}

            <div className="flex gap-3">
              {!isRevealed ? (
                <Button onClick={checkAnswer} disabled={!userInput.trim()}>
                  <Check className="h-4 w-4 mr-2" />
                  Check Answer
                </Button>
              ) : (
                <Button onClick={nextSentence}>
                  {currentSentence < sentences.length - 1 ? "Next Sentence" : "Finish"}
                </Button>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Dictation;