import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { BookOpen, Calendar, Clock, Search, Tag } from "lucide-react";

interface Essay {
  id: string;
  title: string;
  content: string;
  createdAt: string;
  wordCount: number;
  tags?: string[];
}

const TagsPage = () => {
  const { tag } = useParams<{ tag: string }>();
  const navigate = useNavigate();
  const [essays, setEssays] = useState<Essay[]>([]);
  const [filteredEssays, setFilteredEssays] = useState<Essay[]>([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [allTags, setAllTags] = useState<string[]>([]);

  useEffect(() => {
    const essaysData = JSON.parse(localStorage.getItem('essays') || '[]');
    
    // Add mock tags to essays for demo purposes
    const essaysWithTags = essaysData.map((essay: Essay) => ({
      ...essay,
      tags: generateMockTags(essay.title, essay.content)
    }));

    setEssays(essaysWithTags);

    // Extract all unique tags
    const tags = new Set<string>();
    essaysWithTags.forEach((essay: Essay) => {
      essay.tags?.forEach(tag => tags.add(tag));
    });
    setAllTags(Array.from(tags).sort());

    // Filter essays by tag if specified
    if (tag) {
      const filtered = essaysWithTags.filter((essay: Essay) => 
        essay.tags?.includes(decodeURIComponent(tag))
      );
      setFilteredEssays(filtered);
    } else {
      setFilteredEssays(essaysWithTags);
    }
  }, [tag]);

  useEffect(() => {
    if (searchTerm) {
      const filtered = essays.filter(essay =>
        essay.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        essay.content.toLowerCase().includes(searchTerm.toLowerCase()) ||
        essay.tags?.some(tag => tag.toLowerCase().includes(searchTerm.toLowerCase()))
      );
      setFilteredEssays(filtered);
    } else if (tag) {
      const filtered = essays.filter(essay =>
        essay.tags?.includes(decodeURIComponent(tag))
      );
      setFilteredEssays(filtered);
    } else {
      setFilteredEssays(essays);
    }
  }, [searchTerm, essays, tag]);

  const generateMockTags = (title: string, content: string): string[] => {
    const commonTags = ['English', 'Learning', 'Practice'];
    const topicTags = ['Travel', 'Business', 'Technology', 'Culture', 'Education', 'Science'];
    const levelTags = ['Beginner', 'Intermediate', 'Advanced'];
    const typesTags = ['Essay', 'Article', 'Story', 'News'];

    const tags: string[] = [];
    
    // Add common tags
    tags.push(commonTags[Math.floor(Math.random() * commonTags.length)]);
    
    // Add topic based on content
    const lowerContent = (title + ' ' + content).toLowerCase();
    if (lowerContent.includes('travel') || lowerContent.includes('journey')) tags.push('Travel');
    else if (lowerContent.includes('business') || lowerContent.includes('work')) tags.push('Business');
    else if (lowerContent.includes('technology') || lowerContent.includes('computer')) tags.push('Technology');
    else if (lowerContent.includes('culture') || lowerContent.includes('tradition')) tags.push('Culture');
    else if (lowerContent.includes('education') || lowerContent.includes('school')) tags.push('Education');
    else tags.push(topicTags[Math.floor(Math.random() * topicTags.length)]);

    // Add level
    tags.push(levelTags[Math.floor(Math.random() * levelTags.length)]);
    
    // Add type
    tags.push(typesTags[Math.floor(Math.random() * typesTags.length)]);

    return [...new Set(tags)]; // Remove duplicates
  };

  const getReadingTime = (wordCount: number) => {
    return Math.ceil(wordCount / 200);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-foreground mb-2">
          {tag ? `Essays tagged with "${decodeURIComponent(tag)}"` : 'Browse by Tags'}
        </h1>
        <p className="text-muted-foreground">
          {tag ? `Found ${filteredEssays.length} essays` : 'Explore essays organized by topics and difficulty levels'}
        </p>
      </div>

      {/* Search */}
      <Card className="mb-6">
        <CardContent className="pt-6">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
            <Input
              placeholder="Search essays, tags, or content..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10"
            />
          </div>
        </CardContent>
      </Card>

      {/* All Tags Overview */}
      {!tag && (
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Tag className="h-5 w-5" />
              Available Tags
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-2">
              {allTags.map((tagName) => {
                const count = essays.filter(essay => essay.tags?.includes(tagName)).length;
                return (
                  <Badge
                    key={tagName}
                    variant="secondary"
                    className="cursor-pointer hover:bg-primary hover:text-primary-foreground transition-colors"
                    onClick={() => navigate(`/tags/${encodeURIComponent(tagName)}`)}
                  >
                    {tagName} ({count})
                  </Badge>
                );
              })}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Essays Grid */}
      <div className="grid gap-6">
        {filteredEssays.length === 0 ? (
          <Card>
            <CardContent className="text-center py-8">
              <p className="text-muted-foreground mb-4">
                {tag ? `No essays found with tag "${decodeURIComponent(tag)}"` : 'No essays found'}
              </p>
              <Button onClick={() => navigate('/import')}>
                Import Your First Essay
              </Button>
            </CardContent>
          </Card>
        ) : (
          filteredEssays.map((essay) => (
            <Card key={essay.id} className="hover:shadow-md transition-shadow">
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <CardTitle className="text-lg mb-2">{essay.title}</CardTitle>
                    <div className="flex flex-wrap gap-1 mb-3">
                      {essay.tags?.map((tagName) => (
                        <Badge
                          key={tagName}
                          variant="secondary"
                          className="text-xs cursor-pointer hover:bg-primary hover:text-primary-foreground transition-colors"
                          onClick={() => navigate(`/tags/${encodeURIComponent(tagName)}`)}
                        >
                          {tagName}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground mb-4 line-clamp-3">
                  {essay.content.substring(0, 200)}...
                </p>
                
                <div className="flex items-center justify-between text-sm text-muted-foreground mb-4">
                  <div className="flex items-center gap-4">
                    <span className="flex items-center gap-1">
                      <BookOpen className="h-3 w-3" />
                      {essay.wordCount} words
                    </span>
                    <span className="flex items-center gap-1">
                      <Clock className="h-3 w-3" />
                      {getReadingTime(essay.wordCount)} min read
                    </span>
                    <span className="flex items-center gap-1">
                      <Calendar className="h-3 w-3" />
                      {formatDate(essay.createdAt)}
                    </span>
                  </div>
                </div>

                <div className="flex gap-2">
                  <Button onClick={() => navigate(`/essay/${essay.id}`)}>
                    Read Essay
                  </Button>
                  <Button 
                    variant="outline"
                    onClick={() => navigate(`/essay/${essay.id}/recite`)}
                  >
                    Practice
                  </Button>
                  <Button 
                    variant="outline"
                    onClick={() => navigate(`/essay/${essay.id}/cloze`)}
                  >
                    Cloze Test
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>

      {/* Back to all tags */}
      {tag && (
        <div className="mt-8 text-center">
          <Button variant="outline" onClick={() => navigate('/tags')}>
            Browse All Tags
          </Button>
        </div>
      )}
    </div>
  );
};

export default TagsPage;