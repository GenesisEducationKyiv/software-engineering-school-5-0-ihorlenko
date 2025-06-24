const { test, expect } = require('@playwright/test');

test.describe('Weather Subscription', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should load the subscription page', async ({ page }) => {
    await expect(page).toHaveTitle(/Weather Notifier/);
    
    await expect(page.locator('h1')).toContainText('Weather Notifier');
    await expect(page.locator('.subtitle')).toContainText('Get weather updates delivered to your inbox');
    
    await expect(page.locator('#email')).toBeVisible();
    await expect(page.locator('#city')).toBeVisible();
    await expect(page.locator('input[name="frequency"]')).toHaveCount(2);
    await expect(page.locator('.btn-subscribe')).toBeVisible();
  });

  test('should validate required fields', async ({ page }) => {
    await page.click('.btn-subscribe');
    
    const emailInput = page.locator('#email');
    const cityInput = page.locator('#city');
    
    await expect(emailInput).toHaveAttribute('required');
    await expect(cityInput).toHaveAttribute('required');
  });

  test('should validate email format', async ({ page }) => {
    await page.fill('#email', 'invalid-email');
    await page.fill('#city', 'London');
    
    await page.click('.btn-subscribe');
    
    const emailInput = page.locator('#email');
    const isValid = await emailInput.evaluate(el => el.validity.valid);
    expect(isValid).toBeFalsy();
  });

  test('should select frequency options', async ({ page }) => {
    const dailyRadio = page.locator('input[value="daily"]');
    await expect(dailyRadio).toBeChecked();
    
    const hourlyRadio = page.locator('input[value="hourly"]');
    await expect(hourlyRadio).not.toBeChecked();
    
    await hourlyRadio.check();
    await expect(hourlyRadio).toBeChecked();
    await expect(dailyRadio).not.toBeChecked();
    
    await dailyRadio.check();
    await expect(dailyRadio).toBeChecked();
    await expect(hourlyRadio).not.toBeChecked();
  });

  test('should submit valid subscription form', async ({ page }) => {
    await page.route('**/api/weather*', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          city: 'London',
          temperature: 15.5,
          humidity: 65,
          description: 'Partly cloudy'
        })
      });
    });

    await page.route('**/api/subscribe', async route => {
      const request = route.request();
      if (request.method() === 'POST') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            message: 'Please check your email to confirm the subscription'
          })
        });
      } else {
        await route.continue();
      }
    });

    await page.fill('#email', 'test@example.com');
    await page.fill('#city', 'London');
    await page.check('input[value="daily"]');
    
    await page.click('.btn-subscribe');
    
    await expect(page.locator('.message.success')).toBeVisible();
    await expect(page.locator('.message.success')).toContainText('check your email');
    
    await expect(page.locator('#email')).toHaveValue('');
    await expect(page.locator('#city')).toHaveValue('');
  });

  test('should handle API errors gracefully', async ({ page }) => {
    await page.route('**/api/subscribe', async route => {
      await route.fulfill({
        status: 400,
        contentType: 'application/json',
        body: JSON.stringify({
          error: 'Invalid city or weather service unavailable'
        })
      });
    });

    await page.fill('#email', 'test@example.com');
    await page.fill('#city', 'InvalidCity');
    
    await page.click('.btn-subscribe');
    
    await expect(page.locator('.message.error')).toBeVisible();
    await expect(page.locator('.message.error')).toContainText('Invalid city');
  });

  test('should show loading state during submission', async ({ page }) => {
    await page.route('**/api/subscribe', async route => {
      await new Promise(resolve => setTimeout(resolve, 1000));
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          message: 'Please check your email to confirm the subscription'
        })
      });
    });

    await page.fill('#email', 'test@example.com');
    await page.fill('#city', 'London');
    
    await page.click('.btn-subscribe');
    
    const submitButton = page.locator('.btn-subscribe');
    await expect(submitButton).toContainText('Processing...');
    await expect(submitButton).toBeDisabled();
    
    await expect(submitButton).toContainText('Subscribe');
    await expect(submitButton).toBeEnabled();
  });

  test('should handle successful confirmation redirect', async ({ page }) => {
    await page.goto('/?message_type=success&message=Your%20subscription%20has%20been%20successfully%20confirmed!');
    
    await expect(page.locator('.message.success')).toBeVisible();
    await expect(page.locator('.message.success')).toContainText('successfully confirmed');
  });

  test('should handle failed confirmation redirect', async ({ page }) => {
    await page.goto('/?message_type=error&message=Subscription%20not%20found');
    
    await expect(page.locator('.message.error')).toBeVisible();
    await expect(page.locator('.message.error')).toContainText('not found');
  });

  test('should handle successful unsubscribe redirect', async ({ page }) => {
    await page.goto('/?message_type=success&message=You%20have%20successfully%20unsubscribed');
    
    await expect(page.locator('.message.success')).toBeVisible();
    await expect(page.locator('.message.success')).toContainText('unsubscribed');
  });

  test('should be responsive on mobile devices', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    
    await expect(page.locator('#email')).toBeVisible();
    await expect(page.locator('#city')).toBeVisible();
    await expect(page.locator('.btn-subscribe')).toBeVisible();
    
    const radioOptions = page.locator('.radio-options');
    const boundingBox = await radioOptions.boundingBox();
    expect(boundingBox.height).toBeGreaterThan(100);
  });

  test('should have proper accessibility features', async ({ page }) => {
    await expect(page.locator('label[for="email"]')).toBeVisible();
    await expect(page.locator('label[for="city"]')).toBeVisible();
    
    await expect(page.locator('#email')).toHaveAttribute('type', 'email');
    await expect(page.locator('#email')).toHaveAttribute('required');
    await expect(page.locator('#city')).toHaveAttribute('required');
    
    const frequencyRadios = page.locator('input[name="frequency"]');
    await expect(frequencyRadios).toHaveCount(2);
    
    await expect(page.locator('.btn-subscribe')).toContainText('Subscribe');
  });

  test('should handle network connectivity issues', async ({ page }) => {
    await page.route('**/api/subscribe', async route => {
      await route.abort('failed');
    });

    await page.fill('#email', 'test@example.com');
    await page.fill('#city', 'London');
    
    await page.click('.btn-subscribe');
    
    await expect(page.locator('.message.error')).toBeVisible();
    await expect(page.locator('.message.error')).toContainText('error occurred');
  });
});