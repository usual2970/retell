import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";
import { ArrowLeft, Check, RotateCcw, Eye, EyeOff } from "lucide-react";
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

interface ClozeItem {
  id: number;
  text: string;
  isBlank: boolean;
  correctAnswer?: string;
  userAnswer?: string;
  isCorrect?: boolean;
}

const ClozeTest = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const [essay, setEssay] = useState<Essay | null>(null);
  const [clozeItems, setClozeItems] = useState<ClozeItem[]>([]);
  const [currentBlank, setCurrentBlank] = useState(0);
  const [showAnswers, setShowAnswers] = useState(false);
  const [isCompleted, setIsCompleted] = useState(false);
  const [score, setScore] = useState({ correct: 0, total: 0 });

  useEffect(() => {
    if (id) {
      const essays = JSON.parse(localStorage.getItem('essays') || '[]');
      const foundEssay = essays.find((e: Essay) => e.id === id);
      if (foundEssay) {
        setEssay(foundEssay);
        generateClozeTest(foundEssay.content);
      }
    }
  }, [id]);

  const generateClozeTest = (content: string) => {
    const words = content.split(/(\s+)/);
    const items: ClozeItem[] = [];
    let itemId = 0;
    let blankCount = 0;

    // Create blanks for every 7-10 words (excluding very short words)
    for (let i = 0; i < words.length; i++) {
      const word = words[i];
      const cleanWord = word.replace(/[.,!?;:]/g, '');
      
      // Skip whitespace and very short words
      if (word.trim().length === 0) {
        items.push({
          id: itemId++,
          text: word,
          isBlank: false
        });
        continue;
      }

      // Create blank for longer words at intervals
      const shouldBeBlank = cleanWord.length > 3 && 
                           Math.random() > 0.8 && 
                           blankCount < Math.floor(words.filter(w => w.trim().length > 3).length / 8);

      if (shouldBeBlank) {
        items.push({
          id: itemId++,
          text: '_'.repeat(Math.min(cleanWord.length, 12)),
          isBlank: true,
          correctAnswer: cleanWord,
          userAnswer: ''
        });
        blankCount++;
      } else {
        items.push({
          id: itemId++,
          text: word,
          isBlank: false
        });
      }
    }

    setClozeItems(items);
  };

  const handleAnswerChange = (blankIndex: number, value: string) => {
    setClozeItems(prev => 
      prev.map((item, index) => {
        if (item.isBlank && getBlanksBeforeIndex(index) === blankIndex) {
          return { ...item, userAnswer: value };
        }
        return item;
      })
    );
  };

  const getBlanksBeforeIndex = (index: number): number => {
    return clozeItems.slice(0, index).filter(item => item.isBlank).length;
  };

  const checkAnswers = () => {
    let correct = 0;
    let total = 0;

    const updatedItems = clozeItems.map(item => {
      if (item.isBlank) {
        total++;
        const isCorrect = item.userAnswer?.toLowerCase().trim() === 
                         item.correctAnswer?.toLowerCase().trim();
        if (isCorrect) correct++;
        return { ...item, isCorrect };
      }
      return item;
    });

    setClozeItems(updatedItems);
    setScore({ correct, total });
    setIsCompleted(true);
    setShowAnswers(true);

    toast({
      title: "Test Complete!",
      description: `Score: ${correct}/${total} (${Math.round((correct/total) * 100)}%)`,
    });
  };

  const resetTest = () => {
    setClozeItems(prev => 
      prev.map(item => ({
        ...item,
        userAnswer: item.isBlank ? '' : item.userAnswer,
        isCorrect: undefined
      }))
    );
    setCurrentBlank(0);
    setShowAnswers(false);
    setIsCompleted(false);
    setScore({ correct: 0, total: 0 });
  };

  const toggleAnswers = () => {
    setShowAnswers(!showAnswers);
  };

  if (!essay) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p className="text-center text-muted-foreground">Essay not found</p>
      </div>
    );
  }

  const blanks = clozeItems.filter(item => item.isBlank);
  const completedBlanks = blanks.filter(blank => blank.userAnswer?.trim()).length;
  const progress = blanks.length > 0 ? (completedBlanks / blanks.length) * 100 : 0;

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
        <h1 className="text-3xl font-bold text-foreground mb-2">Cloze Test</h1>
        <p className="text-muted-foreground">{essay.title}</p>
      </div>

      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Progress</CardTitle>
              <div className="text-sm text-muted-foreground">
                {completedBlanks}/{blanks.length} completed
              </div>
            </div>
            <Progress value={progress} className="mt-2" />
          </CardHeader>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Fill in the Blanks</CardTitle>
              <div className="flex gap-2">
                {isCompleted && (
                  <Button variant="outline" onClick={toggleAnswers} size="sm">
                    {showAnswers ? <EyeOff className="h-4 w-4 mr-2" /> : <Eye className="h-4 w-4 mr-2" />}
                    {showAnswers ? "Hide" : "Show"} Answers
                  </Button>
                )}
                <Button onClick={resetTest} variant="outline" size="sm">
                  <RotateCcw className="h-4 w-4 mr-2" />
                  Reset
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="prose max-w-none">
              <div className="text-foreground leading-relaxed text-base">
                {clozeItems.map((item, index) => {
                  if (!item.isBlank) {
                    return <span key={item.id}>{item.text}</span>;
                  }

                  const blankIndex = getBlanksBeforeIndex(index);
                  const inputValue = item.userAnswer || '';
                  
                  let inputClass = "inline-block mx-1 px-2 py-1 border rounded text-center min-w-[100px] ";
                  
                  if (showAnswers && item.isCorrect !== undefined) {
                    inputClass += item.isCorrect 
                      ? "border-green-500 bg-green-50 text-green-800" 
                      : "border-red-500 bg-red-50 text-red-800";
                  } else {
                    inputClass += "border-input";
                  }

                  return (
                    <span key={item.id} className="relative">
                      <input
                        type="text"
                        value={inputValue}
                        onChange={(e) => handleAnswerChange(blankIndex, e.target.value)}
                        className={inputClass}
                        placeholder="___"
                        disabled={isCompleted}
                        style={{ width: `${Math.max(100, (item.correctAnswer?.length || 8) * 10)}px` }}
                      />
                      {showAnswers && item.correctAnswer && (
                        <div className="absolute top-full left-1/2 transform -translate-x-1/2 mt-1 px-2 py-1 bg-primary text-primary-foreground text-xs rounded whitespace-nowrap z-10">
                          {item.correctAnswer}
                        </div>
                      )}
                    </span>
                  );
                })}
              </div>
            </div>

            <div className="flex gap-3 mt-6 pt-4 border-t">
              {!isCompleted ? (
                <Button 
                  onClick={checkAnswers} 
                  disabled={completedBlanks === 0}
                  className="flex items-center gap-2"
                >
                  <Check className="h-4 w-4" />
                  Check Answers
                </Button>
              ) : (
                <div className="flex items-center gap-4">
                  <div className="text-lg font-semibold text-foreground">
                    Score: {score.correct}/{score.total} ({Math.round((score.correct/score.total) * 100)}%)
                  </div>
                  <Button onClick={() => navigate(`/essay/${id}/dictation`)} variant="outline">
                    Try Dictation
                  </Button>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default ClozeTest;