function formHandler() {
  return {
    inputText: "",
    shortUrl: "",
    error: "",
    async submitForm() {
      this.error = "";
      this.shortUrl = "";
      
      if (!this.inputText.trim()) {
        this.error = "Please enter a valid URL";
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
      }
    },
    
    async copyToClipboard() {
      try {
        await navigator.clipboard.writeText(this.shortUrl);
        // You could add a temporary success message here
        const button = event.target;
        const originalText = button.textContent;
        button.textContent = 'Copied!';
        button.classList.remove('bg-green-500', 'hover:bg-green-600');
        button.classList.add('bg-green-600');
        
        setTimeout(() => {
          button.textContent = originalText;
          button.classList.remove('bg-green-600');
          button.classList.add('bg-green-500', 'hover:bg-green-600');
        }, 2000);
      } catch (err) {
        // Fallback for older browsers
        const textArea = document.createElement('textarea');
        textArea.value = this.shortUrl;
        document.body.appendChild(textArea);
        textArea.select();
        document.execCommand('copy');
        document.body.removeChild(textArea);
      }
    },
    
    resetForm() {
      this.inputText = "";
      this.shortUrl = "";
      this.error = "";
    }
  };
}
