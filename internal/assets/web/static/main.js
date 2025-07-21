function formHandler() {
  return {
    inputText: "",
    shortUrl: "",
    error: "",
    isLoading: false,
    copied: false,
    async submitForm() {
      this.isLoading = true;
      this.error = "";
      this.shortUrl = "";
      
      if (!this.inputText.trim()) {
        this.error = "Please enter a valid URL";
        this.isLoading = false;
        return;
      }
      
      try {
        const response = await fetch("/create", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ longUrl: this.inputText }),
        });

        const data = await response.json();

        if (!response.ok) {
          throw new Error(data.error || "Failed to create short link");
        }

        // Assuming the API returns the short URL in data.shortUrl or similar
        this.shortUrl = data.shortUrl || data.url || `https://choto.link/${data.id || 'abc123'}`;
      } catch (err) {
        this.error = err.message || "Something went wrong. Please try again.";
      } finally {
        this.isLoading = false;
      }
    },
    
    async copyToClipboard() {
      try {
        await navigator.clipboard.writeText(this.shortUrl);
        this.copied = true;
        
        setTimeout(() => {
          this.copied = false;
        }, 2000);
      } catch (err) {
        // Fallback for older browsers
        const textArea = document.createElement('textarea');
        textArea.value = this.shortUrl;
        document.body.appendChild(textArea);
        textArea.select();
        document.execCommand('copy');
        document.body.removeChild(textArea);
        
        this.copied = true;
        setTimeout(() => {
          this.copied = false;
        }, 1000);
      }
    },
    
    resetForm() {
      this.inputText = "";
      this.shortUrl = "";
      this.error = "";
    }
  };
}
