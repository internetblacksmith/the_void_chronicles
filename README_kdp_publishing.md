# Amazon KDP Publishing Guide 📚

Complete guide for self-publishing "Void Reavers" on Amazon KDP (Kindle Direct Publishing).

## Overview

This guide covers the entire process from manuscript to published book:
1. Creating a print-ready PDF interior
2. Designing a compliant book cover
3. Publishing on Amazon KDP
4. Marketing and distribution tips

## Quick Start

```bash
# Generate your KDP-compliant PDF
./create_kdp_pdf.sh

# Files created:
# - kdp_output/void_reavers_kdp.pdf (interior file)
# - kdp_output/cover_specs.txt (cover dimensions)
```

## Interior File Requirements

### Amazon KDP Specifications

**File Requirements:**
- Format: PDF
- Maximum size: 650 MB
- Minimum resolution: 300 DPI for all images
- Fonts: Must be embedded
- Single pages (not spreads)

**Page Setup:**
- Popular trim size: 6" × 9" (recommended for novels)
- Margins: 0.75" inner, 0.5" outer, 0.75" top/bottom
- Bleed: 0.125" if images extend to page edge

**Content Requirements:**
- No crop marks or printer marks
- No annotations or comments
- Pages must face same direction
- Minimum 24 pages, maximum 828 pages

## Using the KDP PDF Generator

### Step 1: Run the Generator

```bash
./create_kdp_pdf.sh
```

You'll be prompted for:

1. **Trim Size** - The final size of your printed book
   - 5" × 8" - Smaller novel size
   - 6" × 9" - Most common for fiction (recommended)
   - 7" × 10" - Larger format

2. **Page Count** - Affects spine width for cover
   - Script estimates from word count
   - Adjust based on your formatting

3. **Bleed Settings** - Only if you have full-page images
   - Most novels don't need bleed

### Step 2: Review the Generated PDF

Check for:
- ✅ Proper margins and page size
- ✅ Readable fonts (12pt recommended)
- ✅ Consistent formatting throughout
- ✅ Page numbers in correct position
- ✅ Chapter headings properly styled
- ✅ No widows or orphans (single lines)

### Step 3: Create Your Book Cover

Using the specifications in `cover_specs.txt`:

**Free Options:**
- Amazon's Cover Creator (built into KDP)
- Canva (canva.com) - has book cover templates
- BookCoverMaker (bookcovermaker.com)

**Paid Options:**
- 99designs - design contest (~$299+)
- Fiverr - individual designers ($5-200)
- Reedsy - professional designers ($300-800)
- GetCovers - premade covers ($45-300)

**Cover Must Include:**
- Front cover (right side)
- Spine (if 79+ pages)
- Back cover (left side)
- Bleed area (0.125" all sides)

## Publishing Process

### Step 1: Create KDP Account

1. Go to [kdp.amazon.com](https://kdp.amazon.com)
2. Sign in with Amazon account
3. Complete tax information
4. Set up payment method

### Step 2: Create New Paperback

1. Click "Create New Title"
2. Choose "Paperback"
3. Enter book details:
   - Title: "Void Reavers"
   - Subtitle: "A Tale of Space Pirates and Cosmic Plunder"
   - Series: "The Void Chronicles" (Book 1)
   - Author: Your name or pen name
   - Description: Book synopsis
   - Keywords: 7 keywords/phrases for discoverability
   - Categories: Choose 2 (e.g., Science Fiction > Space Opera)

### Step 3: ISBN Options

**Free ISBN from Amazon:**
- ✅ No cost
- ❌ Only for Amazon distribution
- ❌ Amazon listed as publisher

**Your Own ISBN:**
- ✅ Wider distribution options
- ✅ You're listed as publisher
- ❌ Cost: $125 for 1, $295 for 10 (US Bowker)

### Step 4: Upload Content

1. **Interior**: Upload your PDF from `kdp_output/`
2. **Cover**: Upload your cover or use Cover Creator
3. Preview your book using online previewer
4. Download and review the proof PDF

### Step 5: Set Pricing

**Printing Costs** (6×9, 300 pages, black & white):
- US: ~$4.85
- UK: ~£3.80
- EU: ~€4.20

**Royalty Options:**
- 60% royalty (Amazon exclusive territories)
- Minimum list price: Printing cost × 2.5

**Pricing Strategy:**
- Research similar books in your genre
- New authors: $9.99-14.99 paperback
- Established authors: $14.99-19.99

### Step 6: Proof Your Book

**Always Order a Proof Copy!**
- Author copies at printing cost
- Check physical quality
- Verify formatting looks good in print
- Test readability

## Marketing Tips

### Pre-Launch (2-4 weeks before)

1. **Build buzz:**
   - Share cover reveal
   - Post excerpts
   - Create author website
   - Set up Goodreads author profile

2. **Advanced Review Copies (ARCs):**
   - Send to book bloggers
   - NetGalley or BookSirens
   - Goodreads giveaways

### Launch Week

1. **Pricing strategy:**
   - Consider promotional pricing
   - Kindle Countdown Deal

2. **Reviews:**
   - Ask beta readers to review
   - Follow up with ARC readers
   - Join author review exchanges

3. **Promotion:**
   - BookBub Featured Deal (paid)
   - Facebook/Instagram ads
   - Amazon Advertising

### Ongoing Marketing

1. **Series strategy:**
   - Plan Book 2 launch
   - First book as loss leader
   - Box sets later

2. **Email list:**
   - Newsletter signup in book
   - Reader magnet (free story)
   - Regular updates

## Financial Considerations

### Costs to Budget

**Required:**
- ISBN (optional): $125-295
- Proof copy: ~$5-10
- Shipping: ~$5-10

**Recommended:**
- Professional cover: $100-500
- Editing: $500-2000
- Marketing: $100-500

**Total Budget:**
- Minimum: $10 (DIY everything)
- Realistic: $200-500
- Professional: $1000-3000

### Revenue Expectations

**First-Time Authors:**
- Month 1: 10-50 copies
- Month 2-3: 5-20 copies/month
- Ongoing: 1-10 copies/month

**With Marketing:**
- Launch: 50-200 copies
- Monthly: 20-100 copies
- Series boost: 2-5x increase

## Common Issues and Solutions

### PDF Generation Problems

**"LaTeX not found"**
```bash
# Ubuntu/Debian
sudo apt install texlive-full

# macOS
brew install --cask mactex
```

**File too large (>650MB)**
- Reduce image quality
- Remove unnecessary images
- Compress PDF with GhostScript

### KDP Upload Issues

**"Margins outside printable area"**
- Increase margins in LaTeX file
- Standard minimum: 0.5" all sides

**"Fonts not embedded"**
- Regenerate PDF with font embedding
- Use standard fonts (Times, Arial)

**"Page count mismatch"**
- Ensure no blank pages at end
- Check PDF page count matches entered count

## Next Steps

1. **Generate your PDF**: Run `./create_kdp_pdf.sh`
2. **Review carefully**: Check every page
3. **Design cover**: Use specs from `cover_specs.txt`
4. **Set up KDP account**: Complete tax info
5. **Upload and preview**: Use KDP's previewer
6. **Order proof**: Always check physical copy
7. **Launch**: Announce to your audience

## Resources

### Amazon KDP
- [KDP Help](https://kdp.amazon.com/help)
- [Cover Calculator](https://kdp.amazon.com/cover-calculator)
- [Previewer](https://kdp.amazon.com/tools-and-resources)

### Cover Design
- [Canva Book Covers](https://www.canva.com/book-covers/)
- [100 Covers Tutorial](https://100covers.com/)
- [BookBrush](https://bookbrush.com/) - mockups

### Marketing
- [BookFunnel](https://bookfunnel.com/) - reader magnets
- [Written Word Media](https://www.writtenwordmedia.com/) - promotions
- [20BooksTo50K Facebook Group](https://www.facebook.com/groups/20Booksto50k/) - author community

### ISBN Purchase
- US: [Bowker](https://www.myidentifiers.com/)
- UK: [Nielsen](https://www.nielsenisbnstore.com/)
- Canada: [CISS](https://www.bac-lac.gc.ca/eng/services/isbn-canada/)
- Australia: [Thorpe-Bowker](https://www.myidentifiers.com.au/)

---

**Good luck with publishing "Void Reavers"! May your book find its way to readers across the galaxy!** 🚀📚