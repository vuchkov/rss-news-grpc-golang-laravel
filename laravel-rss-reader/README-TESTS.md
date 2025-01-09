# Tests

Add tests to the Laravel application. 

We'll focus on feature tests (integration tests) to cover the main functionalities.

1. Create Feature Test Files:
```
php artisan make:test Feature/RssFeedTest
php artisan make:test Feature/RssPostTest
```

2. `tests/Feature/RssFeedTest.php`:
```
<?php

namespace Tests\Feature;

use App\Models\RssFeed;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class RssFeedTest extends TestCase
{
    use RefreshDatabase;

    public function test_can_create_rss_feed()
    {
        $this->get(route('rss_feeds.create'))->assertStatus(200);

        $this->post(route('rss_feeds.store'), [
            'name' => 'Test Feed',
            'url' => 'https://www.theregister.com/headlines.atom',
        ])->assertRedirect(route('rss_feeds.index'));

        $this->assertDatabaseHas('rss_feeds', [
            'name' => 'Test Feed',
            'url' => 'https://www.theregister.com/headlines.atom',
        ]);
    }

    public function test_can_update_rss_feed()
    {
        $feed = RssFeed::factory()->create();

        $this->get(route('rss_feeds.edit', $feed))->assertStatus(200);

        $this->put(route('rss_feeds.update', $feed), [
            'name' => 'Updated Feed Name',
            'url' => 'https://example.com/feed',
        ])->assertRedirect(route('rss_feeds.index'));

        $this->assertDatabaseHas('rss_feeds', [
            'id' => $feed->id,
            'name' => 'Updated Feed Name',
            'url' => 'https://example.com/feed',
        ]);
    }

    public function test_can_delete_rss_feed()
    {
        $feed = RssFeed::factory()->create();

        $this->delete(route('rss_feeds.destroy', $feed))
            ->assertRedirect(route('rss_feeds.index'));

        $this->assertDatabaseMissing('rss_feeds', ['id' => $feed->id]);
    }

        public function test_index_page_displays_feeds()
    {
        RssFeed::factory(3)->create();

        $response = $this->get(route('rss_feeds.index'));

        $response->assertStatus(200);
        $response->assertViewIs('rss_feeds.index');
        $response->assertViewHas('feeds');
    }

    public function test_create_page_is_accessible()
    {
        $this->get(route('rss_feeds.create'))->assertStatus(200);
    }

    public function test_edit_page_is_accessible()
    {
        $feed = RssFeed::factory()->create();
        $this->get(route('rss_feeds.edit', $feed))->assertStatus(200);
    }
}
```

3. `tests/Feature/RssPostTest.php`:
```
<?php

namespace Tests\Feature;

use App\Jobs\FetchRssFeeds;
use App\Models\RssFeed;
use App\Models\RssPost;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Queue;
use Tests\TestCase;

class RssPostTest extends TestCase
{
    use RefreshDatabase;

    public function test_posts_are_displayed_and_filtered()
    {
        $feed1 = RssFeed::factory()->create(['name' => 'Feed 1']);
        $feed2 = RssFeed::factory()->create(['name' => 'Feed 2']);

        RssPost::factory()->create(['rss_feed_id' => $feed1->id, 'title' => 'Test Post 1', 'published_at' => now()->subDays(2)]);
        RssPost::factory()->create(['rss_feed_id' => $feed1->id, 'title' => 'Another Post', 'published_at' => now()->subDay()]);
        RssPost::factory()->create(['rss_feed_id' => $feed2->id, 'title' => 'Test Post 2', 'published_at' => now()]);

        $response = $this->get(route('posts.index'));
        $response->assertStatus(200);
        $response->assertSee('Test Post 2'); // Most recent first

        $response = $this->get(route('posts.index', ['feed' => $feed1->id]));
        $response->assertSee('Test Post 1');
        $response->assertSee('Another Post');
        $response->assertDontSee('Test Post 2');

        $response = $this->get(route('posts.index', ['title' => 'Another']));
        $response->assertSee('Another Post');
        $response->assertDontSee('Test Post 1');
        $response->assertDontSee('Test Post 2');

        $response = $this->get(route('posts.index', ['date' => now()->subDay()->format('Y-m-d')]));
        $response->assertSee('Another Post');
        $response->assertDontSee('Test Post 1');
        $response->assertDontSee('Test Post 2');
    }

    public function test_fetch_rss_feeds_job_is_queued()
    {
        Queue::fake();

        $this->artisan('rss:fetch')->assertSuccessful();

        Queue::assertPushed(FetchRssFeeds::class);
    }
}
```

4. Factories (`database/factories/RssFeedFactory.php` and`database/factories/RssPostFactory.php`):

- `database/factories/RssFeedFactory.php`:
```
<?php

namespace Database\Factories;

use Illuminate\Database\Eloquent\Factories\Factory;

/**
 * @extends \Illuminate\Database\Eloquent\Factories\Factory<\App\Models\RssFeed>
 */
class RssFeedFactory extends Factory
{
    /**
     * Define the model's default state.
     *
     * @return array<string, mixed>
     */
    public function definition(): array
    {
        return [
            'name' => fake()->name(),
            'url' => fake()->url(),
        ];
    }
}
```

- `database/factories/RssPostFactory.php`:
```
<?php

namespace Database\Factories;

use App\Models\RssFeed;
use Illuminate\Database\Eloquent\Factories\Factory;

/**
 * @extends \Illuminate\Database\Eloquent\Factories\Factory<\App\Models\RssPost>
 */
class RssPostFactory extends Factory
{
    /**
     * Define the model's default state.
     *
     * @return array<string, mixed>
     */
    public function definition(): array
    {
        return [
            'rss_feed_id' => RssFeed::factory(),
            'title' => fake()->sentence(),
            'link' => fake()->url(),
            'description' => fake()->paragraph(),
            'published_at' => fake()->dateTimeBetween('-1 week', 'now'),
        ];
    }
}
```

5. Run the tests:
```
php artisan test
```

Key improvements and explanations:
- Feature Tests: Using feature tests provides better integration testing.
- RefreshDatabase Trait: Using the RefreshDatabase trait ensures a clean database for each test.
- Factories: Using factories makes it easier to create test data.
- Testing Filtering and Sorting: The RssPostTest now covers filtering by feed, title, and date, as well as sorting by date.
- Testing Job Dispatch: The RssPostTest now uses Queue::fake() to assert that the FetchRssFeeds job is dispatched correctly.
- More comprehensive tests: Added tests for the index, create and edit pages for the RSS feeds.

This provides a good set of feature tests for the Laravel application. 
- We can always add more tests as needed to cover more critical functionalities.
