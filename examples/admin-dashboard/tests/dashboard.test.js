// @ts-check
const { test, expect } = require('@playwright/test');

const BASE_URL = 'http://localhost:8081';

test.describe('Admin Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Start the server - assumes it's already running
    await page.goto(BASE_URL);
  });

  test('should load the dashboard page', async ({ page }) => {
    // Check that the page loads
    await expect(page).toHaveTitle('Admin Dashboard');
    
    // Check for main layout elements
    await expect(page.locator('header')).toBeVisible();
    await expect(page.locator('aside')).toBeVisible();
    await expect(page.locator('main')).toBeVisible();
    
    // Check for "Admin Panel" heading in sidebar
    await expect(page.locator('h1', { hasText: 'Admin Panel' })).toBeVisible();
  });

  test('should display metrics grid with correct cards', async ({ page }) => {
    // Check that all metric cards are present
    const metricCards = page.locator('[class*="tw-31"]'); // bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6...
    await expect(metricCards).toHaveCount(4);
    
    // Check specific metric values
    await expect(page.locator('text=Users')).toBeVisible();
    await expect(page.locator('text=12,345')).toBeVisible();
    await expect(page.locator('text=Revenue')).toBeVisible();
    await expect(page.locator('text=$56,789')).toBeVisible();
    await expect(page.locator('text=Orders')).toBeVisible();
    await expect(page.locator('text=1,234')).toBeVisible();
    await expect(page.locator('text=Conversion')).toBeVisible();
    await expect(page.locator('text=3.2%')).toBeVisible();
  });

  test('should display navigation items correctly', async ({ page }) => {
    // Check that navigation items are present
    await expect(page.locator('a[href="/"]', { hasText: 'Dashboard' })).toBeVisible();
    await expect(page.locator('a[href="/users"]', { hasText: 'Users' })).toBeVisible();
    await expect(page.locator('a[href="/products"]', { hasText: 'Products' })).toBeVisible();
    await expect(page.locator('a[href="/orders"]', { hasText: 'Orders' })).toBeVisible();
    await expect(page.locator('a[href="/analytics"]', { hasText: 'Analytics' })).toBeVisible();
    await expect(page.locator('a[href="/settings"]', { hasText: 'Settings' })).toBeVisible();
    
    // Check that Dashboard is marked as active
    const dashboardLink = page.locator('a[href="/"]', { hasText: 'Dashboard' });
    await expect(dashboardLink).toHaveClass(/tw-10/); // Active state class
  });

  test('should display data tables with headers and sample data', async ({ page }) => {
    // Check for "Recent Orders" table
    await expect(page.locator('h2', { hasText: 'Recent Orders' })).toBeVisible();
    
    // Check table headers for Recent Orders
    const ordersTable = page.locator('table').first();
    await expect(ordersTable.locator('th', { hasText: 'Order ID' })).toBeVisible();
    await expect(ordersTable.locator('th', { hasText: 'Customer' })).toBeVisible();
    await expect(ordersTable.locator('th', { hasText: 'Amount' })).toBeVisible();
    await expect(ordersTable.locator('th', { hasText: 'Status' })).toBeVisible();
    
    // Check for "Top Products" table
    await expect(page.locator('h2', { hasText: 'Top Products' })).toBeVisible();
    
    // Check table headers for Top Products
    const productsTable = page.locator('table').nth(1);
    await expect(productsTable.locator('th', { hasText: 'Product' })).toBeVisible();
    await expect(productsTable.locator('th', { hasText: 'Sales' })).toBeVisible();
    await expect(productsTable.locator('th', { hasText: 'Revenue' })).toBeVisible();
    await expect(productsTable.locator('th', { hasText: 'Stock' })).toBeVisible();
  });

  test('should show status badges in tables', async ({ page }) => {
    // Check for Active/Inactive status badges
    await expect(page.locator('span', { hasText: 'Active' })).toBeVisible();
    await expect(page.locator('span', { hasText: 'Inactive' })).toBeVisible();
    
    // Check that status badges have correct styling
    const activeBadge = page.locator('span', { hasText: 'Active' }).first();
    await expect(activeBadge).toHaveClass(/tw-57/); // Green styling class
    
    const inactiveBadge = page.locator('span', { hasText: 'Inactive' }).first();
    await expect(inactiveBadge).toHaveClass(/tw-60/); // Red styling class
  });

  test('should have working search input', async ({ page }) => {
    const searchInput = page.locator('input[type="search"]');
    await expect(searchInput).toBeVisible();
    await expect(searchInput).toHaveAttribute('placeholder', 'Search...');
    
    // Test that we can type in the search input
    await searchInput.click();
    await searchInput.fill('test search');
    await expect(searchInput).toHaveValue('test search');
  });

  test('should display header with profile section', async ({ page }) => {
    // Check for profile image
    await expect(page.locator('img[alt="Profile"]')).toBeVisible();
    
    // Check for profile name
    await expect(page.locator('span', { hasText: 'John Doe' })).toBeVisible();
    
    // Check for notification bell (with red dot)
    const notificationButton = page.locator('button').filter({ has: page.locator('svg') }).nth(1);
    await expect(notificationButton).toBeVisible();
    await expect(notificationButton.locator('.tw-25')).toBeVisible(); // Red notification dot
  });

  test('should test sidebar responsiveness with collapsed state', async ({ page }) => {
    // Test with sidebar collapsed query parameter
    await page.goto(`${BASE_URL}?sidebar=collapsed`);
    
    // The sidebar should still be visible but with collapsed transform
    const sidebar = page.locator('aside');
    await expect(sidebar).toBeVisible();
    
    // On mobile, check that menu button is visible
    await page.setViewportSize({ width: 640, height: 480 });
    const menuButton = page.locator('button').first();
    await expect(menuButton).toBeVisible();
  });

  test('should have proper table sorting indicators', async ({ page }) => {
    // Check for sortable headers with sort indicators
    const amountHeader = page.locator('th', { hasText: 'Amount' });
    await expect(amountHeader).toHaveClass(/tw-48/); // Sorted state styling
    
    const revenueHeader = page.locator('th', { hasText: 'Revenue' });
    await expect(revenueHeader).toHaveClass(/tw-48/); // Sorted state styling
    
    // Check for sort arrow icons
    await expect(amountHeader.locator('svg')).toBeVisible();
    await expect(revenueHeader.locator('svg')).toBeVisible();
  });

  test('should validate twerge class generation', async ({ page }) => {
    // Check that twerge-generated classes are being applied
    const bodyElement = page.locator('body');
    await expect(bodyElement).toHaveClass(/tw-0/); // min-h-screen bg-gray-50 dark:bg-gray-900
    
    // Check main container
    const mainContainer = page.locator('div').first();
    await expect(mainContainer).toHaveClass(/tw-1/); // flex h-screen overflow-hidden
    
    // Check sidebar
    const sidebar = page.locator('aside');
    await expect(sidebar).toHaveClass(/tw-4/); // Complex sidebar classes
    
    // Check header
    const header = page.locator('header');
    await expect(header).toHaveClass(/tw-16/); // Header styling
  });

  test('should test navigation hover states', async ({ page }) => {
    // Test navigation item hover effects
    const usersLink = page.locator('a[href="/users"]');
    
    // Hover over the Users link
    await usersLink.hover();
    
    // Check that it has the submenu indicator (group-hover effect should be working)
    const chevron = usersLink.locator('svg').last();
    await expect(chevron).toBeVisible();
  });

  test('should validate metric cards show correct trends', async ({ page }) => {
    // Check that positive trends are green
    const positiveChange = page.locator('span', { hasText: '↗ 12%' });
    await expect(positiveChange).toHaveClass(/tw-36/); // Green text styling
    
    // Check that negative trends are red  
    const negativeChange = page.locator('span', { hasText: '↘ 3%' });
    await expect(negativeChange).toHaveClass(/tw-38/); // Red text styling
  });

  test('should test accessibility features', async ({ page }) => {
    // Check that main landmarks are properly labeled
    await expect(page.locator('main')).toBeVisible();
    await expect(page.locator('header')).toBeVisible();
    await expect(page.locator('aside')).toBeVisible();
    
    // Check that form elements have proper labels
    const searchInput = page.locator('input[type="search"]');
    await expect(searchInput).toHaveAttribute('placeholder', 'Search...');
    
    // Check that buttons are properly accessible
    const buttons = page.locator('button');
    const buttonCount = await buttons.count();
    expect(buttonCount).toBeGreaterThan(0);
    
    // Check that images have alt text
    const profileImage = page.locator('img[alt="Profile"]');
    await expect(profileImage).toBeVisible();
  });

  test('should test table interactivity', async ({ page }) => {
    // Check that table rows have hover effects
    const tableRows = page.locator('tbody tr');
    const firstRow = tableRows.first();
    
    await firstRow.hover();
    await expect(firstRow).toHaveClass(/tw-53/); // hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors duration-200
  });
});

test.describe('Dashboard Error Handling', () => {
  test('should handle server not running', async ({ page }) => {
    // Test what happens when trying to connect to a non-running server
    try {
      await page.goto('http://localhost:9999', { timeout: 5000 });
    } catch (error) {
      // This is expected when server is not running
      expect(error.message).toContain('net::ERR_CONNECTION_REFUSED');
    }
  });
});

test.describe('Dashboard Responsive Design', () => {
  test('should work on tablet view', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.goto(BASE_URL);
    
    // Check that layout adapts appropriately
    await expect(page.locator('aside')).toBeVisible();
    await expect(page.locator('main')).toBeVisible();
    
    // Metrics grid should stack differently
    const metricsGrid = page.locator('[class*="tw-30"]');
    await expect(metricsGrid).toBeVisible();
  });

  test('should work on mobile view', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto(BASE_URL);
    
    // Mobile menu button should be visible
    const menuButton = page.locator('button').first();
    await expect(menuButton).toBeVisible();
    
    // Check that content is still accessible
    await expect(page.locator('main')).toBeVisible();
  });
});