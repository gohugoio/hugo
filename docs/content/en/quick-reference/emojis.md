---
title: Emojis
description: Include emoji shortcodes in your Markdown or templates.
categories: [quick reference]
keywords: [emoji]
menu:
  docs:
    parent: quick-reference
    weight: 20
weight: 20
toc: true
---

## Attribution

This quick reference guide was generated using the [ikatyang/emoji-cheat-sheet] project which reads from the [GitHub Emoji API] and the [Unicode Full Emoji List].

Note that GitHub [custom emoji] are not supported.

[custom emoji]: #github-custom-emoji
[github emoji api]: https://api.github.com/emojis
[ikatyang/emoji-cheat-sheet]: https://github.com/ikatyang/emoji-cheat-sheet/
[unicode full emoji list]: https://unicode.org/emoji/charts/full-emoji-list.html

## Usage

Configure Hugo to enable emoji processing in Markdown:

{{< code-toggle file=hugo >}}
enableEmoji = true
{{< /code-toggle >}}

With emoji processing enabled, this Markdown:

```md
Hello! :wave:
```

Is rendered to:

```html
Hello! &#x1f44b;
```

And in your browser... Hello! :wave:

To process an emoji shortcode from within a template, use the [`emojify`] function or pass the string through the [`RenderString`] method on a `Page` object:

```go-html-template
{{ "Hello! :wave:" | .RenderString }}
```

[`emojify`]: /functions/transform/emojify/
[`RenderString`]: /methods/page/renderstring/

<!-- 
To generate the sections below:

    git clone https://github.com/ikatyang/emoji-cheat-sheet
    cd emoji-cheat-sheet
    npm install
    npm run generate

Then...

    1. Copy and paste from README.md
    2. Search/replace (regex) "^###\s" with "## "
    3. Search/replace "^####\s " with "### "
-->

## Table of Contents

- [Smileys & Emotion](#smileys--emotion)
- [People & Body](#people--body)
- [Animals & Nature](#animals--nature)
- [Food & Drink](#food--drink)
- [Travel & Places](#travel--places)
- [Activities](#activities)
- [Objects](#objects)
- [Symbols](#symbols)
- [Flags](#flags)
- [GitHub Custom Emoji](#github-custom-emoji)

## Smileys & Emotion

- [Face Smiling](#face-smiling)
- [Face Affection](#face-affection)
- [Face Tongue](#face-tongue)
- [Face Hand](#face-hand)
- [Face Neutral Skeptical](#face-neutral-skeptical)
- [Face Sleepy](#face-sleepy)
- [Face Unwell](#face-unwell)
- [Face Hat](#face-hat)
- [Face Glasses](#face-glasses)
- [Face Concerned](#face-concerned)
- [Face Negative](#face-negative)
- [Face Costume](#face-costume)
- [Cat Face](#cat-face)
- [Monkey Face](#monkey-face)
- [Heart](#heart)
- [Emotion](#emotion)

### Face Smiling

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :grinning: | `:grinning:` | :smiley: | `:smiley:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :smile: | `:smile:` | :grin: | `:grin:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :laughing: | `:laughing:` <br /> `:satisfied:` | :sweat_smile: | `:sweat_smile:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :rofl: | `:rofl:` | :joy: | `:joy:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :slightly_smiling_face: | `:slightly_smiling_face:` | :upside_down_face: | `:upside_down_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :melting_face: | `:melting_face:` | :wink: | `:wink:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :blush: | `:blush:` | :innocent: | `:innocent:` | [top](#table-of-contents) |

### Face Affection

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :smiling_face_with_three_hearts: | `:smiling_face_with_three_hearts:` | :heart_eyes: | `:heart_eyes:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :star_struck: | `:star_struck:` | :kissing_heart: | `:kissing_heart:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :kissing: | `:kissing:` | :relaxed: | `:relaxed:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :kissing_closed_eyes: | `:kissing_closed_eyes:` | :kissing_smiling_eyes: | `:kissing_smiling_eyes:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :smiling_face_with_tear: | `:smiling_face_with_tear:` | | | [top](#table-of-contents) |

### Face Tongue

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :yum: | `:yum:` | :stuck_out_tongue: | `:stuck_out_tongue:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :stuck_out_tongue_winking_eye: | `:stuck_out_tongue_winking_eye:` | :zany_face: | `:zany_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :stuck_out_tongue_closed_eyes: | `:stuck_out_tongue_closed_eyes:` | :money_mouth_face: | `:money_mouth_face:` | [top](#table-of-contents) |

### Face Hand

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :hugs: | `:hugs:` | :hand_over_mouth: | `:hand_over_mouth:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :face_with_open_eyes_and_hand_over_mouth: | `:face_with_open_eyes_and_hand_over_mouth:` | :face_with_peeking_eye: | `:face_with_peeking_eye:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :shushing_face: | `:shushing_face:` | :thinking: | `:thinking:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :saluting_face: | `:saluting_face:` | | | [top](#table-of-contents) |

### Face Neutral Skeptical

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :zipper_mouth_face: | `:zipper_mouth_face:` | :raised_eyebrow: | `:raised_eyebrow:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :neutral_face: | `:neutral_face:` | :expressionless: | `:expressionless:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :no_mouth: | `:no_mouth:` | :dotted_line_face: | `:dotted_line_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :face_in_clouds: | `:face_in_clouds:` | :smirk: | `:smirk:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :unamused: | `:unamused:` | :roll_eyes: | `:roll_eyes:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :grimacing: | `:grimacing:` | :face_exhaling: | `:face_exhaling:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :lying_face: | `:lying_face:` | :shaking_face: | `:shaking_face:` | [top](#table-of-contents) |

### Face Sleepy

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :relieved: | `:relieved:` | :pensive: | `:pensive:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :sleepy: | `:sleepy:` | :drooling_face: | `:drooling_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :sleeping: | `:sleeping:` | | | [top](#table-of-contents) |

### Face Unwell

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :mask: | `:mask:` | :face_with_thermometer: | `:face_with_thermometer:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :face_with_head_bandage: | `:face_with_head_bandage:` | :nauseated_face: | `:nauseated_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :vomiting_face: | `:vomiting_face:` | :sneezing_face: | `:sneezing_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :hot_face: | `:hot_face:` | :cold_face: | `:cold_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :woozy_face: | `:woozy_face:` | :dizzy_face: | `:dizzy_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :face_with_spiral_eyes: | `:face_with_spiral_eyes:` | :exploding_head: | `:exploding_head:` | [top](#table-of-contents) |

### Face Hat

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :cowboy_hat_face: | `:cowboy_hat_face:` | :partying_face: | `:partying_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :disguised_face: | `:disguised_face:` | | | [top](#table-of-contents) |

### Face Glasses

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :sunglasses: | `:sunglasses:` | :nerd_face: | `:nerd_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :monocle_face: | `:monocle_face:` | | | [top](#table-of-contents) |

### Face Concerned

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :confused: | `:confused:` | :face_with_diagonal_mouth: | `:face_with_diagonal_mouth:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :worried: | `:worried:` | :slightly_frowning_face: | `:slightly_frowning_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :frowning_face: | `:frowning_face:` | :open_mouth: | `:open_mouth:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :hushed: | `:hushed:` | :astonished: | `:astonished:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :flushed: | `:flushed:` | :pleading_face: | `:pleading_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :face_holding_back_tears: | `:face_holding_back_tears:` | :frowning: | `:frowning:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :anguished: | `:anguished:` | :fearful: | `:fearful:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :cold_sweat: | `:cold_sweat:` | :disappointed_relieved: | `:disappointed_relieved:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :cry: | `:cry:` | :sob: | `:sob:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :scream: | `:scream:` | :confounded: | `:confounded:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :persevere: | `:persevere:` | :disappointed: | `:disappointed:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :sweat: | `:sweat:` | :weary: | `:weary:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :tired_face: | `:tired_face:` | :yawning_face: | `:yawning_face:` | [top](#table-of-contents) |

### Face Negative

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :triumph: | `:triumph:` | :pout: | `:pout:` <br /> `:rage:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :angry: | `:angry:` | :cursing_face: | `:cursing_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :smiling_imp: | `:smiling_imp:` | :imp: | `:imp:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :skull: | `:skull:` | :skull_and_crossbones: | `:skull_and_crossbones:` | [top](#table-of-contents) |

### Face Costume

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :hankey: | `:hankey:` <br /> `:poop:` <br /> `:shit:` | :clown_face: | `:clown_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :japanese_ogre: | `:japanese_ogre:` | :japanese_goblin: | `:japanese_goblin:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :ghost: | `:ghost:` | :alien: | `:alien:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :space_invader: | `:space_invader:` | :robot: | `:robot:` | [top](#table-of-contents) |

### Cat Face

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :smiley_cat: | `:smiley_cat:` | :smile_cat: | `:smile_cat:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :joy_cat: | `:joy_cat:` | :heart_eyes_cat: | `:heart_eyes_cat:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :smirk_cat: | `:smirk_cat:` | :kissing_cat: | `:kissing_cat:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :scream_cat: | `:scream_cat:` | :crying_cat_face: | `:crying_cat_face:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :pouting_cat: | `:pouting_cat:` | | | [top](#table-of-contents) |

### Monkey Face

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :see_no_evil: | `:see_no_evil:` | :hear_no_evil: | `:hear_no_evil:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :speak_no_evil: | `:speak_no_evil:` | | | [top](#table-of-contents) |

### Heart

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :love_letter: | `:love_letter:` | :cupid: | `:cupid:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :gift_heart: | `:gift_heart:` | :sparkling_heart: | `:sparkling_heart:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :heartpulse: | `:heartpulse:` | :heartbeat: | `:heartbeat:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :revolving_hearts: | `:revolving_hearts:` | :two_hearts: | `:two_hearts:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :heart_decoration: | `:heart_decoration:` | :heavy_heart_exclamation: | `:heavy_heart_exclamation:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :broken_heart: | `:broken_heart:` | :heart_on_fire: | `:heart_on_fire:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :mending_heart: | `:mending_heart:` | :heart: | `:heart:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :pink_heart: | `:pink_heart:` | :orange_heart: | `:orange_heart:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :yellow_heart: | `:yellow_heart:` | :green_heart: | `:green_heart:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :blue_heart: | `:blue_heart:` | :light_blue_heart: | `:light_blue_heart:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :purple_heart: | `:purple_heart:` | :brown_heart: | `:brown_heart:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :black_heart: | `:black_heart:` | :grey_heart: | `:grey_heart:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :white_heart: | `:white_heart:` | | | [top](#table-of-contents) |

### Emotion

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#smileys--emotion) | :kiss: | `:kiss:` | :100: | `:100:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :anger: | `:anger:` | :boom: | `:boom:` <br /> `:collision:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :dizzy: | `:dizzy:` | :sweat_drops: | `:sweat_drops:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :dash: | `:dash:` | :hole: | `:hole:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :speech_balloon: | `:speech_balloon:` | :eye_speech_bubble: | `:eye_speech_bubble:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :left_speech_bubble: | `:left_speech_bubble:` | :right_anger_bubble: | `:right_anger_bubble:` | [top](#table-of-contents) |
| [top](#smileys--emotion) | :thought_balloon: | `:thought_balloon:` | :zzz: | `:zzz:` | [top](#table-of-contents) |

## People & Body

- [Hand Fingers Open](#hand-fingers-open)
- [Hand Fingers Partial](#hand-fingers-partial)
- [Hand Single Finger](#hand-single-finger)
- [Hand Fingers Closed](#hand-fingers-closed)
- [Hands](#hands)
- [Hand Prop](#hand-prop)
- [Body Parts](#body-parts)
- [Person](#person)
- [Person Gesture](#person-gesture)
- [Person Role](#person-role)
- [Person Fantasy](#person-fantasy)
- [Person Activity](#person-activity)
- [Person Sport](#person-sport)
- [Person Resting](#person-resting)
- [Family](#family)
- [Person Symbol](#person-symbol)

### Hand Fingers Open

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :wave: | `:wave:` | :raised_back_of_hand: | `:raised_back_of_hand:` | [top](#table-of-contents) |
| [top](#people--body) | :raised_hand_with_fingers_splayed: | `:raised_hand_with_fingers_splayed:` | :hand: | `:hand:` <br /> `:raised_hand:` | [top](#table-of-contents) |
| [top](#people--body) | :vulcan_salute: | `:vulcan_salute:` | :rightwards_hand: | `:rightwards_hand:` | [top](#table-of-contents) |
| [top](#people--body) | :leftwards_hand: | `:leftwards_hand:` | :palm_down_hand: | `:palm_down_hand:` | [top](#table-of-contents) |
| [top](#people--body) | :palm_up_hand: | `:palm_up_hand:` | :leftwards_pushing_hand: | `:leftwards_pushing_hand:` | [top](#table-of-contents) |
| [top](#people--body) | :rightwards_pushing_hand: | `:rightwards_pushing_hand:` | | | [top](#table-of-contents) |

### Hand Fingers Partial

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :ok_hand: | `:ok_hand:` | :pinched_fingers: | `:pinched_fingers:` | [top](#table-of-contents) |
| [top](#people--body) | :pinching_hand: | `:pinching_hand:` | :v: | `:v:` | [top](#table-of-contents) |
| [top](#people--body) | :crossed_fingers: | `:crossed_fingers:` | :hand_with_index_finger_and_thumb_crossed: | `:hand_with_index_finger_and_thumb_crossed:` | [top](#table-of-contents) |
| [top](#people--body) | :love_you_gesture: | `:love_you_gesture:` | :metal: | `:metal:` | [top](#table-of-contents) |
| [top](#people--body) | :call_me_hand: | `:call_me_hand:` | | | [top](#table-of-contents) |

### Hand Single Finger

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :point_left: | `:point_left:` | :point_right: | `:point_right:` | [top](#table-of-contents) |
| [top](#people--body) | :point_up_2: | `:point_up_2:` | :fu: | `:fu:` <br /> `:middle_finger:` | [top](#table-of-contents) |
| [top](#people--body) | :point_down: | `:point_down:` | :point_up: | `:point_up:` | [top](#table-of-contents) |
| [top](#people--body) | :index_pointing_at_the_viewer: | `:index_pointing_at_the_viewer:` | | | [top](#table-of-contents) |

### Hand Fingers Closed

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :+1: | `:+1:` <br /> `:thumbsup:` | :-1: | `:-1:` <br /> `:thumbsdown:` | [top](#table-of-contents) |
| [top](#people--body) | :fist: | `:fist:` <br /> `:fist_raised:` | :facepunch: | `:facepunch:` <br /> `:fist_oncoming:` <br /> `:punch:` | [top](#table-of-contents) |
| [top](#people--body) | :fist_left: | `:fist_left:` | :fist_right: | `:fist_right:` | [top](#table-of-contents) |

### Hands

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :clap: | `:clap:` | :raised_hands: | `:raised_hands:` | [top](#table-of-contents) |
| [top](#people--body) | :heart_hands: | `:heart_hands:` | :open_hands: | `:open_hands:` | [top](#table-of-contents) |
| [top](#people--body) | :palms_up_together: | `:palms_up_together:` | :handshake: | `:handshake:` | [top](#table-of-contents) |
| [top](#people--body) | :pray: | `:pray:` | | | [top](#table-of-contents) |

### Hand Prop

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :writing_hand: | `:writing_hand:` | :nail_care: | `:nail_care:` | [top](#table-of-contents) |
| [top](#people--body) | :selfie: | `:selfie:` | | | [top](#table-of-contents) |

### Body Parts

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :muscle: | `:muscle:` | :mechanical_arm: | `:mechanical_arm:` | [top](#table-of-contents) |
| [top](#people--body) | :mechanical_leg: | `:mechanical_leg:` | :leg: | `:leg:` | [top](#table-of-contents) |
| [top](#people--body) | :foot: | `:foot:` | :ear: | `:ear:` | [top](#table-of-contents) |
| [top](#people--body) | :ear_with_hearing_aid: | `:ear_with_hearing_aid:` | :nose: | `:nose:` | [top](#table-of-contents) |
| [top](#people--body) | :brain: | `:brain:` | :anatomical_heart: | `:anatomical_heart:` | [top](#table-of-contents) |
| [top](#people--body) | :lungs: | `:lungs:` | :tooth: | `:tooth:` | [top](#table-of-contents) |
| [top](#people--body) | :bone: | `:bone:` | :eyes: | `:eyes:` | [top](#table-of-contents) |
| [top](#people--body) | :eye: | `:eye:` | :tongue: | `:tongue:` | [top](#table-of-contents) |
| [top](#people--body) | :lips: | `:lips:` | :biting_lip: | `:biting_lip:` | [top](#table-of-contents) |

### Person

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :baby: | `:baby:` | :child: | `:child:` | [top](#table-of-contents) |
| [top](#people--body) | :boy: | `:boy:` | :girl: | `:girl:` | [top](#table-of-contents) |
| [top](#people--body) | :adult: | `:adult:` | :blond_haired_person: | `:blond_haired_person:` | [top](#table-of-contents) |
| [top](#people--body) | :man: | `:man:` | :bearded_person: | `:bearded_person:` | [top](#table-of-contents) |
| [top](#people--body) | :man_beard: | `:man_beard:` | :woman_beard: | `:woman_beard:` | [top](#table-of-contents) |
| [top](#people--body) | :red_haired_man: | `:red_haired_man:` | :curly_haired_man: | `:curly_haired_man:` | [top](#table-of-contents) |
| [top](#people--body) | :white_haired_man: | `:white_haired_man:` | :bald_man: | `:bald_man:` | [top](#table-of-contents) |
| [top](#people--body) | :woman: | `:woman:` | :red_haired_woman: | `:red_haired_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :person_red_hair: | `:person_red_hair:` | :curly_haired_woman: | `:curly_haired_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :person_curly_hair: | `:person_curly_hair:` | :white_haired_woman: | `:white_haired_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :person_white_hair: | `:person_white_hair:` | :bald_woman: | `:bald_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :person_bald: | `:person_bald:` | :blond_haired_woman: | `:blond_haired_woman:` <br /> `:blonde_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :blond_haired_man: | `:blond_haired_man:` | :older_adult: | `:older_adult:` | [top](#table-of-contents) |
| [top](#people--body) | :older_man: | `:older_man:` | :older_woman: | `:older_woman:` | [top](#table-of-contents) |

### Person Gesture

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :frowning_person: | `:frowning_person:` | :frowning_man: | `:frowning_man:` | [top](#table-of-contents) |
| [top](#people--body) | :frowning_woman: | `:frowning_woman:` | :pouting_face: | `:pouting_face:` | [top](#table-of-contents) |
| [top](#people--body) | :pouting_man: | `:pouting_man:` | :pouting_woman: | `:pouting_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :no_good: | `:no_good:` | :ng_man: | `:ng_man:` <br /> `:no_good_man:` | [top](#table-of-contents) |
| [top](#people--body) | :ng_woman: | `:ng_woman:` <br /> `:no_good_woman:` | :ok_person: | `:ok_person:` | [top](#table-of-contents) |
| [top](#people--body) | :ok_man: | `:ok_man:` | :ok_woman: | `:ok_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :information_desk_person: | `:information_desk_person:` <br /> `:tipping_hand_person:` | :sassy_man: | `:sassy_man:` <br /> `:tipping_hand_man:` | [top](#table-of-contents) |
| [top](#people--body) | :sassy_woman: | `:sassy_woman:` <br /> `:tipping_hand_woman:` | :raising_hand: | `:raising_hand:` | [top](#table-of-contents) |
| [top](#people--body) | :raising_hand_man: | `:raising_hand_man:` | :raising_hand_woman: | `:raising_hand_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :deaf_person: | `:deaf_person:` | :deaf_man: | `:deaf_man:` | [top](#table-of-contents) |
| [top](#people--body) | :deaf_woman: | `:deaf_woman:` | :bow: | `:bow:` | [top](#table-of-contents) |
| [top](#people--body) | :bowing_man: | `:bowing_man:` | :bowing_woman: | `:bowing_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :facepalm: | `:facepalm:` | :man_facepalming: | `:man_facepalming:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_facepalming: | `:woman_facepalming:` | :shrug: | `:shrug:` | [top](#table-of-contents) |
| [top](#people--body) | :man_shrugging: | `:man_shrugging:` | :woman_shrugging: | `:woman_shrugging:` | [top](#table-of-contents) |

### Person Role

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :health_worker: | `:health_worker:` | :man_health_worker: | `:man_health_worker:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_health_worker: | `:woman_health_worker:` | :student: | `:student:` | [top](#table-of-contents) |
| [top](#people--body) | :man_student: | `:man_student:` | :woman_student: | `:woman_student:` | [top](#table-of-contents) |
| [top](#people--body) | :teacher: | `:teacher:` | :man_teacher: | `:man_teacher:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_teacher: | `:woman_teacher:` | :judge: | `:judge:` | [top](#table-of-contents) |
| [top](#people--body) | :man_judge: | `:man_judge:` | :woman_judge: | `:woman_judge:` | [top](#table-of-contents) |
| [top](#people--body) | :farmer: | `:farmer:` | :man_farmer: | `:man_farmer:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_farmer: | `:woman_farmer:` | :cook: | `:cook:` | [top](#table-of-contents) |
| [top](#people--body) | :man_cook: | `:man_cook:` | :woman_cook: | `:woman_cook:` | [top](#table-of-contents) |
| [top](#people--body) | :mechanic: | `:mechanic:` | :man_mechanic: | `:man_mechanic:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_mechanic: | `:woman_mechanic:` | :factory_worker: | `:factory_worker:` | [top](#table-of-contents) |
| [top](#people--body) | :man_factory_worker: | `:man_factory_worker:` | :woman_factory_worker: | `:woman_factory_worker:` | [top](#table-of-contents) |
| [top](#people--body) | :office_worker: | `:office_worker:` | :man_office_worker: | `:man_office_worker:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_office_worker: | `:woman_office_worker:` | :scientist: | `:scientist:` | [top](#table-of-contents) |
| [top](#people--body) | :man_scientist: | `:man_scientist:` | :woman_scientist: | `:woman_scientist:` | [top](#table-of-contents) |
| [top](#people--body) | :technologist: | `:technologist:` | :man_technologist: | `:man_technologist:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_technologist: | `:woman_technologist:` | :singer: | `:singer:` | [top](#table-of-contents) |
| [top](#people--body) | :man_singer: | `:man_singer:` | :woman_singer: | `:woman_singer:` | [top](#table-of-contents) |
| [top](#people--body) | :artist: | `:artist:` | :man_artist: | `:man_artist:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_artist: | `:woman_artist:` | :pilot: | `:pilot:` | [top](#table-of-contents) |
| [top](#people--body) | :man_pilot: | `:man_pilot:` | :woman_pilot: | `:woman_pilot:` | [top](#table-of-contents) |
| [top](#people--body) | :astronaut: | `:astronaut:` | :man_astronaut: | `:man_astronaut:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_astronaut: | `:woman_astronaut:` | :firefighter: | `:firefighter:` | [top](#table-of-contents) |
| [top](#people--body) | :man_firefighter: | `:man_firefighter:` | :woman_firefighter: | `:woman_firefighter:` | [top](#table-of-contents) |
| [top](#people--body) | :cop: | `:cop:` <br /> `:police_officer:` | :policeman: | `:policeman:` | [top](#table-of-contents) |
| [top](#people--body) | :policewoman: | `:policewoman:` | :detective: | `:detective:` | [top](#table-of-contents) |
| [top](#people--body) | :male_detective: | `:male_detective:` | :female_detective: | `:female_detective:` | [top](#table-of-contents) |
| [top](#people--body) | :guard: | `:guard:` | :guardsman: | `:guardsman:` | [top](#table-of-contents) |
| [top](#people--body) | :guardswoman: | `:guardswoman:` | :ninja: | `:ninja:` | [top](#table-of-contents) |
| [top](#people--body) | :construction_worker: | `:construction_worker:` | :construction_worker_man: | `:construction_worker_man:` | [top](#table-of-contents) |
| [top](#people--body) | :construction_worker_woman: | `:construction_worker_woman:` | :person_with_crown: | `:person_with_crown:` | [top](#table-of-contents) |
| [top](#people--body) | :prince: | `:prince:` | :princess: | `:princess:` | [top](#table-of-contents) |
| [top](#people--body) | :person_with_turban: | `:person_with_turban:` | :man_with_turban: | `:man_with_turban:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_with_turban: | `:woman_with_turban:` | :man_with_gua_pi_mao: | `:man_with_gua_pi_mao:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_with_headscarf: | `:woman_with_headscarf:` | :person_in_tuxedo: | `:person_in_tuxedo:` | [top](#table-of-contents) |
| [top](#people--body) | :man_in_tuxedo: | `:man_in_tuxedo:` | :woman_in_tuxedo: | `:woman_in_tuxedo:` | [top](#table-of-contents) |
| [top](#people--body) | :person_with_veil: | `:person_with_veil:` | :man_with_veil: | `:man_with_veil:` | [top](#table-of-contents) |
| [top](#people--body) | :bride_with_veil: | `:bride_with_veil:` <br /> `:woman_with_veil:` | :pregnant_woman: | `:pregnant_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :pregnant_man: | `:pregnant_man:` | :pregnant_person: | `:pregnant_person:` | [top](#table-of-contents) |
| [top](#people--body) | :breast_feeding: | `:breast_feeding:` | :woman_feeding_baby: | `:woman_feeding_baby:` | [top](#table-of-contents) |
| [top](#people--body) | :man_feeding_baby: | `:man_feeding_baby:` | :person_feeding_baby: | `:person_feeding_baby:` | [top](#table-of-contents) |

### Person Fantasy

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :angel: | `:angel:` | :santa: | `:santa:` | [top](#table-of-contents) |
| [top](#people--body) | :mrs_claus: | `:mrs_claus:` | :mx_claus: | `:mx_claus:` | [top](#table-of-contents) |
| [top](#people--body) | :superhero: | `:superhero:` | :superhero_man: | `:superhero_man:` | [top](#table-of-contents) |
| [top](#people--body) | :superhero_woman: | `:superhero_woman:` | :supervillain: | `:supervillain:` | [top](#table-of-contents) |
| [top](#people--body) | :supervillain_man: | `:supervillain_man:` | :supervillain_woman: | `:supervillain_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :mage: | `:mage:` | :mage_man: | `:mage_man:` | [top](#table-of-contents) |
| [top](#people--body) | :mage_woman: | `:mage_woman:` | :fairy: | `:fairy:` | [top](#table-of-contents) |
| [top](#people--body) | :fairy_man: | `:fairy_man:` | :fairy_woman: | `:fairy_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :vampire: | `:vampire:` | :vampire_man: | `:vampire_man:` | [top](#table-of-contents) |
| [top](#people--body) | :vampire_woman: | `:vampire_woman:` | :merperson: | `:merperson:` | [top](#table-of-contents) |
| [top](#people--body) | :merman: | `:merman:` | :mermaid: | `:mermaid:` | [top](#table-of-contents) |
| [top](#people--body) | :elf: | `:elf:` | :elf_man: | `:elf_man:` | [top](#table-of-contents) |
| [top](#people--body) | :elf_woman: | `:elf_woman:` | :genie: | `:genie:` | [top](#table-of-contents) |
| [top](#people--body) | :genie_man: | `:genie_man:` | :genie_woman: | `:genie_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :zombie: | `:zombie:` | :zombie_man: | `:zombie_man:` | [top](#table-of-contents) |
| [top](#people--body) | :zombie_woman: | `:zombie_woman:` | :troll: | `:troll:` | [top](#table-of-contents) |

### Person Activity

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :massage: | `:massage:` | :massage_man: | `:massage_man:` | [top](#table-of-contents) |
| [top](#people--body) | :massage_woman: | `:massage_woman:` | :haircut: | `:haircut:` | [top](#table-of-contents) |
| [top](#people--body) | :haircut_man: | `:haircut_man:` | :haircut_woman: | `:haircut_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :walking: | `:walking:` | :walking_man: | `:walking_man:` | [top](#table-of-contents) |
| [top](#people--body) | :walking_woman: | `:walking_woman:` | :standing_person: | `:standing_person:` | [top](#table-of-contents) |
| [top](#people--body) | :standing_man: | `:standing_man:` | :standing_woman: | `:standing_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :kneeling_person: | `:kneeling_person:` | :kneeling_man: | `:kneeling_man:` | [top](#table-of-contents) |
| [top](#people--body) | :kneeling_woman: | `:kneeling_woman:` | :person_with_probing_cane: | `:person_with_probing_cane:` | [top](#table-of-contents) |
| [top](#people--body) | :man_with_probing_cane: | `:man_with_probing_cane:` | :woman_with_probing_cane: | `:woman_with_probing_cane:` | [top](#table-of-contents) |
| [top](#people--body) | :person_in_motorized_wheelchair: | `:person_in_motorized_wheelchair:` | :man_in_motorized_wheelchair: | `:man_in_motorized_wheelchair:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_in_motorized_wheelchair: | `:woman_in_motorized_wheelchair:` | :person_in_manual_wheelchair: | `:person_in_manual_wheelchair:` | [top](#table-of-contents) |
| [top](#people--body) | :man_in_manual_wheelchair: | `:man_in_manual_wheelchair:` | :woman_in_manual_wheelchair: | `:woman_in_manual_wheelchair:` | [top](#table-of-contents) |
| [top](#people--body) | :runner: | `:runner:` <br /> `:running:` | :running_man: | `:running_man:` | [top](#table-of-contents) |
| [top](#people--body) | :running_woman: | `:running_woman:` | :dancer: | `:dancer:` <br /> `:woman_dancing:` | [top](#table-of-contents) |
| [top](#people--body) | :man_dancing: | `:man_dancing:` | :business_suit_levitating: | `:business_suit_levitating:` | [top](#table-of-contents) |
| [top](#people--body) | :dancers: | `:dancers:` | :dancing_men: | `:dancing_men:` | [top](#table-of-contents) |
| [top](#people--body) | :dancing_women: | `:dancing_women:` | :sauna_person: | `:sauna_person:` | [top](#table-of-contents) |
| [top](#people--body) | :sauna_man: | `:sauna_man:` | :sauna_woman: | `:sauna_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :climbing: | `:climbing:` | :climbing_man: | `:climbing_man:` | [top](#table-of-contents) |
| [top](#people--body) | :climbing_woman: | `:climbing_woman:` | | | [top](#table-of-contents) |

### Person Sport

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :person_fencing: | `:person_fencing:` | :horse_racing: | `:horse_racing:` | [top](#table-of-contents) |
| [top](#people--body) | :skier: | `:skier:` | :snowboarder: | `:snowboarder:` | [top](#table-of-contents) |
| [top](#people--body) | :golfing: | `:golfing:` | :golfing_man: | `:golfing_man:` | [top](#table-of-contents) |
| [top](#people--body) | :golfing_woman: | `:golfing_woman:` | :surfer: | `:surfer:` | [top](#table-of-contents) |
| [top](#people--body) | :surfing_man: | `:surfing_man:` | :surfing_woman: | `:surfing_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :rowboat: | `:rowboat:` | :rowing_man: | `:rowing_man:` | [top](#table-of-contents) |
| [top](#people--body) | :rowing_woman: | `:rowing_woman:` | :swimmer: | `:swimmer:` | [top](#table-of-contents) |
| [top](#people--body) | :swimming_man: | `:swimming_man:` | :swimming_woman: | `:swimming_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :bouncing_ball_person: | `:bouncing_ball_person:` | :basketball_man: | `:basketball_man:` <br /> `:bouncing_ball_man:` | [top](#table-of-contents) |
| [top](#people--body) | :basketball_woman: | `:basketball_woman:` <br /> `:bouncing_ball_woman:` | :weight_lifting: | `:weight_lifting:` | [top](#table-of-contents) |
| [top](#people--body) | :weight_lifting_man: | `:weight_lifting_man:` | :weight_lifting_woman: | `:weight_lifting_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :bicyclist: | `:bicyclist:` | :biking_man: | `:biking_man:` | [top](#table-of-contents) |
| [top](#people--body) | :biking_woman: | `:biking_woman:` | :mountain_bicyclist: | `:mountain_bicyclist:` | [top](#table-of-contents) |
| [top](#people--body) | :mountain_biking_man: | `:mountain_biking_man:` | :mountain_biking_woman: | `:mountain_biking_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :cartwheeling: | `:cartwheeling:` | :man_cartwheeling: | `:man_cartwheeling:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_cartwheeling: | `:woman_cartwheeling:` | :wrestling: | `:wrestling:` | [top](#table-of-contents) |
| [top](#people--body) | :men_wrestling: | `:men_wrestling:` | :women_wrestling: | `:women_wrestling:` | [top](#table-of-contents) |
| [top](#people--body) | :water_polo: | `:water_polo:` | :man_playing_water_polo: | `:man_playing_water_polo:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_playing_water_polo: | `:woman_playing_water_polo:` | :handball_person: | `:handball_person:` | [top](#table-of-contents) |
| [top](#people--body) | :man_playing_handball: | `:man_playing_handball:` | :woman_playing_handball: | `:woman_playing_handball:` | [top](#table-of-contents) |
| [top](#people--body) | :juggling_person: | `:juggling_person:` | :man_juggling: | `:man_juggling:` | [top](#table-of-contents) |
| [top](#people--body) | :woman_juggling: | `:woman_juggling:` | | | [top](#table-of-contents) |

### Person Resting

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :lotus_position: | `:lotus_position:` | :lotus_position_man: | `:lotus_position_man:` | [top](#table-of-contents) |
| [top](#people--body) | :lotus_position_woman: | `:lotus_position_woman:` | :bath: | `:bath:` | [top](#table-of-contents) |
| [top](#people--body) | :sleeping_bed: | `:sleeping_bed:` | | | [top](#table-of-contents) |

### Family

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :people_holding_hands: | `:people_holding_hands:` | :two_women_holding_hands: | `:two_women_holding_hands:` | [top](#table-of-contents) |
| [top](#people--body) | :couple: | `:couple:` | :two_men_holding_hands: | `:two_men_holding_hands:` | [top](#table-of-contents) |
| [top](#people--body) | :couplekiss: | `:couplekiss:` | :couplekiss_man_woman: | `:couplekiss_man_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :couplekiss_man_man: | `:couplekiss_man_man:` | :couplekiss_woman_woman: | `:couplekiss_woman_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :couple_with_heart: | `:couple_with_heart:` | :couple_with_heart_woman_man: | `:couple_with_heart_woman_man:` | [top](#table-of-contents) |
| [top](#people--body) | :couple_with_heart_man_man: | `:couple_with_heart_man_man:` | :couple_with_heart_woman_woman: | `:couple_with_heart_woman_woman:` | [top](#table-of-contents) |
| [top](#people--body) | :family_man_woman_boy: | `:family_man_woman_boy:` | :family_man_woman_girl: | `:family_man_woman_girl:` | [top](#table-of-contents) |
| [top](#people--body) | :family_man_woman_girl_boy: | `:family_man_woman_girl_boy:` | :family_man_woman_boy_boy: | `:family_man_woman_boy_boy:` | [top](#table-of-contents) |
| [top](#people--body) | :family_man_woman_girl_girl: | `:family_man_woman_girl_girl:` | :family_man_man_boy: | `:family_man_man_boy:` | [top](#table-of-contents) |
| [top](#people--body) | :family_man_man_girl: | `:family_man_man_girl:` | :family_man_man_girl_boy: | `:family_man_man_girl_boy:` | [top](#table-of-contents) |
| [top](#people--body) | :family_man_man_boy_boy: | `:family_man_man_boy_boy:` | :family_man_man_girl_girl: | `:family_man_man_girl_girl:` | [top](#table-of-contents) |
| [top](#people--body) | :family_woman_woman_boy: | `:family_woman_woman_boy:` | :family_woman_woman_girl: | `:family_woman_woman_girl:` | [top](#table-of-contents) |
| [top](#people--body) | :family_woman_woman_girl_boy: | `:family_woman_woman_girl_boy:` | :family_woman_woman_boy_boy: | `:family_woman_woman_boy_boy:` | [top](#table-of-contents) |
| [top](#people--body) | :family_woman_woman_girl_girl: | `:family_woman_woman_girl_girl:` | :family_man_boy: | `:family_man_boy:` | [top](#table-of-contents) |
| [top](#people--body) | :family_man_boy_boy: | `:family_man_boy_boy:` | :family_man_girl: | `:family_man_girl:` | [top](#table-of-contents) |
| [top](#people--body) | :family_man_girl_boy: | `:family_man_girl_boy:` | :family_man_girl_girl: | `:family_man_girl_girl:` | [top](#table-of-contents) |
| [top](#people--body) | :family_woman_boy: | `:family_woman_boy:` | :family_woman_boy_boy: | `:family_woman_boy_boy:` | [top](#table-of-contents) |
| [top](#people--body) | :family_woman_girl: | `:family_woman_girl:` | :family_woman_girl_boy: | `:family_woman_girl_boy:` | [top](#table-of-contents) |
| [top](#people--body) | :family_woman_girl_girl: | `:family_woman_girl_girl:` | | | [top](#table-of-contents) |

### Person Symbol

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#people--body) | :speaking_head: | `:speaking_head:` | :bust_in_silhouette: | `:bust_in_silhouette:` | [top](#table-of-contents) |
| [top](#people--body) | :busts_in_silhouette: | `:busts_in_silhouette:` | :people_hugging: | `:people_hugging:` | [top](#table-of-contents) |
| [top](#people--body) | :family: | `:family:` | :footprints: | `:footprints:` | [top](#table-of-contents) |

## Animals & Nature

- [Animal Mammal](#animal-mammal)
- [Animal Bird](#animal-bird)
- [Animal Amphibian](#animal-amphibian)
- [Animal Reptile](#animal-reptile)
- [Animal Marine](#animal-marine)
- [Animal Bug](#animal-bug)
- [Plant Flower](#plant-flower)
- [Plant Other](#plant-other)

### Animal Mammal

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#animals--nature) | :monkey_face: | `:monkey_face:` | :monkey: | `:monkey:` | [top](#table-of-contents) |
| [top](#animals--nature) | :gorilla: | `:gorilla:` | :orangutan: | `:orangutan:` | [top](#table-of-contents) |
| [top](#animals--nature) | :dog: | `:dog:` | :dog2: | `:dog2:` | [top](#table-of-contents) |
| [top](#animals--nature) | :guide_dog: | `:guide_dog:` | :service_dog: | `:service_dog:` | [top](#table-of-contents) |
| [top](#animals--nature) | :poodle: | `:poodle:` | :wolf: | `:wolf:` | [top](#table-of-contents) |
| [top](#animals--nature) | :fox_face: | `:fox_face:` | :raccoon: | `:raccoon:` | [top](#table-of-contents) |
| [top](#animals--nature) | :cat: | `:cat:` | :cat2: | `:cat2:` | [top](#table-of-contents) |
| [top](#animals--nature) | :black_cat: | `:black_cat:` | :lion: | `:lion:` | [top](#table-of-contents) |
| [top](#animals--nature) | :tiger: | `:tiger:` | :tiger2: | `:tiger2:` | [top](#table-of-contents) |
| [top](#animals--nature) | :leopard: | `:leopard:` | :horse: | `:horse:` | [top](#table-of-contents) |
| [top](#animals--nature) | :moose: | `:moose:` | :donkey: | `:donkey:` | [top](#table-of-contents) |
| [top](#animals--nature) | :racehorse: | `:racehorse:` | :unicorn: | `:unicorn:` | [top](#table-of-contents) |
| [top](#animals--nature) | :zebra: | `:zebra:` | :deer: | `:deer:` | [top](#table-of-contents) |
| [top](#animals--nature) | :bison: | `:bison:` | :cow: | `:cow:` | [top](#table-of-contents) |
| [top](#animals--nature) | :ox: | `:ox:` | :water_buffalo: | `:water_buffalo:` | [top](#table-of-contents) |
| [top](#animals--nature) | :cow2: | `:cow2:` | :pig: | `:pig:` | [top](#table-of-contents) |
| [top](#animals--nature) | :pig2: | `:pig2:` | :boar: | `:boar:` | [top](#table-of-contents) |
| [top](#animals--nature) | :pig_nose: | `:pig_nose:` | :ram: | `:ram:` | [top](#table-of-contents) |
| [top](#animals--nature) | :sheep: | `:sheep:` | :goat: | `:goat:` | [top](#table-of-contents) |
| [top](#animals--nature) | :dromedary_camel: | `:dromedary_camel:` | :camel: | `:camel:` | [top](#table-of-contents) |
| [top](#animals--nature) | :llama: | `:llama:` | :giraffe: | `:giraffe:` | [top](#table-of-contents) |
| [top](#animals--nature) | :elephant: | `:elephant:` | :mammoth: | `:mammoth:` | [top](#table-of-contents) |
| [top](#animals--nature) | :rhinoceros: | `:rhinoceros:` | :hippopotamus: | `:hippopotamus:` | [top](#table-of-contents) |
| [top](#animals--nature) | :mouse: | `:mouse:` | :mouse2: | `:mouse2:` | [top](#table-of-contents) |
| [top](#animals--nature) | :rat: | `:rat:` | :hamster: | `:hamster:` | [top](#table-of-contents) |
| [top](#animals--nature) | :rabbit: | `:rabbit:` | :rabbit2: | `:rabbit2:` | [top](#table-of-contents) |
| [top](#animals--nature) | :chipmunk: | `:chipmunk:` | :beaver: | `:beaver:` | [top](#table-of-contents) |
| [top](#animals--nature) | :hedgehog: | `:hedgehog:` | :bat: | `:bat:` | [top](#table-of-contents) |
| [top](#animals--nature) | :bear: | `:bear:` | :polar_bear: | `:polar_bear:` | [top](#table-of-contents) |
| [top](#animals--nature) | :koala: | `:koala:` | :panda_face: | `:panda_face:` | [top](#table-of-contents) |
| [top](#animals--nature) | :sloth: | `:sloth:` | :otter: | `:otter:` | [top](#table-of-contents) |
| [top](#animals--nature) | :skunk: | `:skunk:` | :kangaroo: | `:kangaroo:` | [top](#table-of-contents) |
| [top](#animals--nature) | :badger: | `:badger:` | :feet: | `:feet:` <br /> `:paw_prints:` | [top](#table-of-contents) |

### Animal Bird

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#animals--nature) | :turkey: | `:turkey:` | :chicken: | `:chicken:` | [top](#table-of-contents) |
| [top](#animals--nature) | :rooster: | `:rooster:` | :hatching_chick: | `:hatching_chick:` | [top](#table-of-contents) |
| [top](#animals--nature) | :baby_chick: | `:baby_chick:` | :hatched_chick: | `:hatched_chick:` | [top](#table-of-contents) |
| [top](#animals--nature) | :bird: | `:bird:` | :penguin: | `:penguin:` | [top](#table-of-contents) |
| [top](#animals--nature) | :dove: | `:dove:` | :eagle: | `:eagle:` | [top](#table-of-contents) |
| [top](#animals--nature) | :duck: | `:duck:` | :swan: | `:swan:` | [top](#table-of-contents) |
| [top](#animals--nature) | :owl: | `:owl:` | :dodo: | `:dodo:` | [top](#table-of-contents) |
| [top](#animals--nature) | :feather: | `:feather:` | :flamingo: | `:flamingo:` | [top](#table-of-contents) |
| [top](#animals--nature) | :peacock: | `:peacock:` | :parrot: | `:parrot:` | [top](#table-of-contents) |
| [top](#animals--nature) | :wing: | `:wing:` | :black_bird: | `:black_bird:` | [top](#table-of-contents) |
| [top](#animals--nature) | :goose: | `:goose:` | | | [top](#table-of-contents) |

### Animal Amphibian

| | ico | shortcode | |
| - | :-: | - | - |
| [top](#animals--nature) | :frog: | `:frog:` | [top](#table-of-contents) |

### Animal Reptile

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#animals--nature) | :crocodile: | `:crocodile:` | :turtle: | `:turtle:` | [top](#table-of-contents) |
| [top](#animals--nature) | :lizard: | `:lizard:` | :snake: | `:snake:` | [top](#table-of-contents) |
| [top](#animals--nature) | :dragon_face: | `:dragon_face:` | :dragon: | `:dragon:` | [top](#table-of-contents) |
| [top](#animals--nature) | :sauropod: | `:sauropod:` | :t-rex: | `:t-rex:` | [top](#table-of-contents) |

### Animal Marine

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#animals--nature) | :whale: | `:whale:` | :whale2: | `:whale2:` | [top](#table-of-contents) |
| [top](#animals--nature) | :dolphin: | `:dolphin:` <br /> `:flipper:` | :seal: | `:seal:` | [top](#table-of-contents) |
| [top](#animals--nature) | :fish: | `:fish:` | :tropical_fish: | `:tropical_fish:` | [top](#table-of-contents) |
| [top](#animals--nature) | :blowfish: | `:blowfish:` | :shark: | `:shark:` | [top](#table-of-contents) |
| [top](#animals--nature) | :octopus: | `:octopus:` | :shell: | `:shell:` | [top](#table-of-contents) |
| [top](#animals--nature) | :coral: | `:coral:` | :jellyfish: | `:jellyfish:` | [top](#table-of-contents) |

### Animal Bug

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#animals--nature) | :snail: | `:snail:` | :butterfly: | `:butterfly:` | [top](#table-of-contents) |
| [top](#animals--nature) | :bug: | `:bug:` | :ant: | `:ant:` | [top](#table-of-contents) |
| [top](#animals--nature) | :bee: | `:bee:` <br /> `:honeybee:` | :beetle: | `:beetle:` | [top](#table-of-contents) |
| [top](#animals--nature) | :lady_beetle: | `:lady_beetle:` | :cricket: | `:cricket:` | [top](#table-of-contents) |
| [top](#animals--nature) | :cockroach: | `:cockroach:` | :spider: | `:spider:` | [top](#table-of-contents) |
| [top](#animals--nature) | :spider_web: | `:spider_web:` | :scorpion: | `:scorpion:` | [top](#table-of-contents) |
| [top](#animals--nature) | :mosquito: | `:mosquito:` | :fly: | `:fly:` | [top](#table-of-contents) |
| [top](#animals--nature) | :worm: | `:worm:` | :microbe: | `:microbe:` | [top](#table-of-contents) |

### Plant Flower

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#animals--nature) | :bouquet: | `:bouquet:` | :cherry_blossom: | `:cherry_blossom:` | [top](#table-of-contents) |
| [top](#animals--nature) | :white_flower: | `:white_flower:` | :lotus: | `:lotus:` | [top](#table-of-contents) |
| [top](#animals--nature) | :rosette: | `:rosette:` | :rose: | `:rose:` | [top](#table-of-contents) |
| [top](#animals--nature) | :wilted_flower: | `:wilted_flower:` | :hibiscus: | `:hibiscus:` | [top](#table-of-contents) |
| [top](#animals--nature) | :sunflower: | `:sunflower:` | :blossom: | `:blossom:` | [top](#table-of-contents) |
| [top](#animals--nature) | :tulip: | `:tulip:` | :hyacinth: | `:hyacinth:` | [top](#table-of-contents) |

### Plant Other

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#animals--nature) | :seedling: | `:seedling:` | :potted_plant: | `:potted_plant:` | [top](#table-of-contents) |
| [top](#animals--nature) | :evergreen_tree: | `:evergreen_tree:` | :deciduous_tree: | `:deciduous_tree:` | [top](#table-of-contents) |
| [top](#animals--nature) | :palm_tree: | `:palm_tree:` | :cactus: | `:cactus:` | [top](#table-of-contents) |
| [top](#animals--nature) | :ear_of_rice: | `:ear_of_rice:` | :herb: | `:herb:` | [top](#table-of-contents) |
| [top](#animals--nature) | :shamrock: | `:shamrock:` | :four_leaf_clover: | `:four_leaf_clover:` | [top](#table-of-contents) |
| [top](#animals--nature) | :maple_leaf: | `:maple_leaf:` | :fallen_leaf: | `:fallen_leaf:` | [top](#table-of-contents) |
| [top](#animals--nature) | :leaves: | `:leaves:` | :empty_nest: | `:empty_nest:` | [top](#table-of-contents) |
| [top](#animals--nature) | :nest_with_eggs: | `:nest_with_eggs:` | :mushroom: | `:mushroom:` | [top](#table-of-contents) |

## Food & Drink

- [Food Fruit](#food-fruit)
- [Food Vegetable](#food-vegetable)
- [Food Prepared](#food-prepared)
- [Food Asian](#food-asian)
- [Food Marine](#food-marine)
- [Food Sweet](#food-sweet)
- [Drink](#drink)
- [Dishware](#dishware)

### Food Fruit

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#food--drink) | :grapes: | `:grapes:` | :melon: | `:melon:` | [top](#table-of-contents) |
| [top](#food--drink) | :watermelon: | `:watermelon:` | :mandarin: | `:mandarin:` <br /> `:orange:` <br /> `:tangerine:` | [top](#table-of-contents) |
| [top](#food--drink) | :lemon: | `:lemon:` | :banana: | `:banana:` | [top](#table-of-contents) |
| [top](#food--drink) | :pineapple: | `:pineapple:` | :mango: | `:mango:` | [top](#table-of-contents) |
| [top](#food--drink) | :apple: | `:apple:` | :green_apple: | `:green_apple:` | [top](#table-of-contents) |
| [top](#food--drink) | :pear: | `:pear:` | :peach: | `:peach:` | [top](#table-of-contents) |
| [top](#food--drink) | :cherries: | `:cherries:` | :strawberry: | `:strawberry:` | [top](#table-of-contents) |
| [top](#food--drink) | :blueberries: | `:blueberries:` | :kiwi_fruit: | `:kiwi_fruit:` | [top](#table-of-contents) |
| [top](#food--drink) | :tomato: | `:tomato:` | :olive: | `:olive:` | [top](#table-of-contents) |
| [top](#food--drink) | :coconut: | `:coconut:` | | | [top](#table-of-contents) |

### Food Vegetable

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#food--drink) | :avocado: | `:avocado:` | :eggplant: | `:eggplant:` | [top](#table-of-contents) |
| [top](#food--drink) | :potato: | `:potato:` | :carrot: | `:carrot:` | [top](#table-of-contents) |
| [top](#food--drink) | :corn: | `:corn:` | :hot_pepper: | `:hot_pepper:` | [top](#table-of-contents) |
| [top](#food--drink) | :bell_pepper: | `:bell_pepper:` | :cucumber: | `:cucumber:` | [top](#table-of-contents) |
| [top](#food--drink) | :leafy_green: | `:leafy_green:` | :broccoli: | `:broccoli:` | [top](#table-of-contents) |
| [top](#food--drink) | :garlic: | `:garlic:` | :onion: | `:onion:` | [top](#table-of-contents) |
| [top](#food--drink) | :peanuts: | `:peanuts:` | :beans: | `:beans:` | [top](#table-of-contents) |
| [top](#food--drink) | :chestnut: | `:chestnut:` | :ginger_root: | `:ginger_root:` | [top](#table-of-contents) |
| [top](#food--drink) | :pea_pod: | `:pea_pod:` | | | [top](#table-of-contents) |

### Food Prepared

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#food--drink) | :bread: | `:bread:` | :croissant: | `:croissant:` | [top](#table-of-contents) |
| [top](#food--drink) | :baguette_bread: | `:baguette_bread:` | :flatbread: | `:flatbread:` | [top](#table-of-contents) |
| [top](#food--drink) | :pretzel: | `:pretzel:` | :bagel: | `:bagel:` | [top](#table-of-contents) |
| [top](#food--drink) | :pancakes: | `:pancakes:` | :waffle: | `:waffle:` | [top](#table-of-contents) |
| [top](#food--drink) | :cheese: | `:cheese:` | :meat_on_bone: | `:meat_on_bone:` | [top](#table-of-contents) |
| [top](#food--drink) | :poultry_leg: | `:poultry_leg:` | :cut_of_meat: | `:cut_of_meat:` | [top](#table-of-contents) |
| [top](#food--drink) | :bacon: | `:bacon:` | :hamburger: | `:hamburger:` | [top](#table-of-contents) |
| [top](#food--drink) | :fries: | `:fries:` | :pizza: | `:pizza:` | [top](#table-of-contents) |
| [top](#food--drink) | :hotdog: | `:hotdog:` | :sandwich: | `:sandwich:` | [top](#table-of-contents) |
| [top](#food--drink) | :taco: | `:taco:` | :burrito: | `:burrito:` | [top](#table-of-contents) |
| [top](#food--drink) | :tamale: | `:tamale:` | :stuffed_flatbread: | `:stuffed_flatbread:` | [top](#table-of-contents) |
| [top](#food--drink) | :falafel: | `:falafel:` | :egg: | `:egg:` | [top](#table-of-contents) |
| [top](#food--drink) | :fried_egg: | `:fried_egg:` | :shallow_pan_of_food: | `:shallow_pan_of_food:` | [top](#table-of-contents) |
| [top](#food--drink) | :stew: | `:stew:` | :fondue: | `:fondue:` | [top](#table-of-contents) |
| [top](#food--drink) | :bowl_with_spoon: | `:bowl_with_spoon:` | :green_salad: | `:green_salad:` | [top](#table-of-contents) |
| [top](#food--drink) | :popcorn: | `:popcorn:` | :butter: | `:butter:` | [top](#table-of-contents) |
| [top](#food--drink) | :salt: | `:salt:` | :canned_food: | `:canned_food:` | [top](#table-of-contents) |

### Food Asian

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#food--drink) | :bento: | `:bento:` | :rice_cracker: | `:rice_cracker:` | [top](#table-of-contents) |
| [top](#food--drink) | :rice_ball: | `:rice_ball:` | :rice: | `:rice:` | [top](#table-of-contents) |
| [top](#food--drink) | :curry: | `:curry:` | :ramen: | `:ramen:` | [top](#table-of-contents) |
| [top](#food--drink) | :spaghetti: | `:spaghetti:` | :sweet_potato: | `:sweet_potato:` | [top](#table-of-contents) |
| [top](#food--drink) | :oden: | `:oden:` | :sushi: | `:sushi:` | [top](#table-of-contents) |
| [top](#food--drink) | :fried_shrimp: | `:fried_shrimp:` | :fish_cake: | `:fish_cake:` | [top](#table-of-contents) |
| [top](#food--drink) | :moon_cake: | `:moon_cake:` | :dango: | `:dango:` | [top](#table-of-contents) |
| [top](#food--drink) | :dumpling: | `:dumpling:` | :fortune_cookie: | `:fortune_cookie:` | [top](#table-of-contents) |
| [top](#food--drink) | :takeout_box: | `:takeout_box:` | | | [top](#table-of-contents) |

### Food Marine

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#food--drink) | :crab: | `:crab:` | :lobster: | `:lobster:` | [top](#table-of-contents) |
| [top](#food--drink) | :shrimp: | `:shrimp:` | :squid: | `:squid:` | [top](#table-of-contents) |
| [top](#food--drink) | :oyster: | `:oyster:` | | | [top](#table-of-contents) |

### Food Sweet

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#food--drink) | :icecream: | `:icecream:` | :shaved_ice: | `:shaved_ice:` | [top](#table-of-contents) |
| [top](#food--drink) | :ice_cream: | `:ice_cream:` | :doughnut: | `:doughnut:` | [top](#table-of-contents) |
| [top](#food--drink) | :cookie: | `:cookie:` | :birthday: | `:birthday:` | [top](#table-of-contents) |
| [top](#food--drink) | :cake: | `:cake:` | :cupcake: | `:cupcake:` | [top](#table-of-contents) |
| [top](#food--drink) | :pie: | `:pie:` | :chocolate_bar: | `:chocolate_bar:` | [top](#table-of-contents) |
| [top](#food--drink) | :candy: | `:candy:` | :lollipop: | `:lollipop:` | [top](#table-of-contents) |
| [top](#food--drink) | :custard: | `:custard:` | :honey_pot: | `:honey_pot:` | [top](#table-of-contents) |

### Drink

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#food--drink) | :baby_bottle: | `:baby_bottle:` | :milk_glass: | `:milk_glass:` | [top](#table-of-contents) |
| [top](#food--drink) | :coffee: | `:coffee:` | :teapot: | `:teapot:` | [top](#table-of-contents) |
| [top](#food--drink) | :tea: | `:tea:` | :sake: | `:sake:` | [top](#table-of-contents) |
| [top](#food--drink) | :champagne: | `:champagne:` | :wine_glass: | `:wine_glass:` | [top](#table-of-contents) |
| [top](#food--drink) | :cocktail: | `:cocktail:` | :tropical_drink: | `:tropical_drink:` | [top](#table-of-contents) |
| [top](#food--drink) | :beer: | `:beer:` | :beers: | `:beers:` | [top](#table-of-contents) |
| [top](#food--drink) | :clinking_glasses: | `:clinking_glasses:` | :tumbler_glass: | `:tumbler_glass:` | [top](#table-of-contents) |
| [top](#food--drink) | :pouring_liquid: | `:pouring_liquid:` | :cup_with_straw: | `:cup_with_straw:` | [top](#table-of-contents) |
| [top](#food--drink) | :bubble_tea: | `:bubble_tea:` | :beverage_box: | `:beverage_box:` | [top](#table-of-contents) |
| [top](#food--drink) | :mate: | `:mate:` | :ice_cube: | `:ice_cube:` | [top](#table-of-contents) |

### Dishware

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#food--drink) | :chopsticks: | `:chopsticks:` | :plate_with_cutlery: | `:plate_with_cutlery:` | [top](#table-of-contents) |
| [top](#food--drink) | :fork_and_knife: | `:fork_and_knife:` | :spoon: | `:spoon:` | [top](#table-of-contents) |
| [top](#food--drink) | :hocho: | `:hocho:` <br /> `:knife:` | :jar: | `:jar:` | [top](#table-of-contents) |
| [top](#food--drink) | :amphora: | `:amphora:` | | | [top](#table-of-contents) |

## Travel & Places

- [Place Map](#place-map)
- [Place Geographic](#place-geographic)
- [Place Building](#place-building)
- [Place Religious](#place-religious)
- [Place Other](#place-other)
- [Transport Ground](#transport-ground)
- [Transport Water](#transport-water)
- [Transport Air](#transport-air)
- [Hotel](#hotel)
- [Time](#time)
- [Sky & Weather](#sky--weather)

### Place Map

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :earth_africa: | `:earth_africa:` | :earth_americas: | `:earth_americas:` | [top](#table-of-contents) |
| [top](#travel--places) | :earth_asia: | `:earth_asia:` | :globe_with_meridians: | `:globe_with_meridians:` | [top](#table-of-contents) |
| [top](#travel--places) | :world_map: | `:world_map:` | :japan: | `:japan:` | [top](#table-of-contents) |
| [top](#travel--places) | :compass: | `:compass:` | | | [top](#table-of-contents) |

### Place Geographic

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :mountain_snow: | `:mountain_snow:` | :mountain: | `:mountain:` | [top](#table-of-contents) |
| [top](#travel--places) | :volcano: | `:volcano:` | :mount_fuji: | `:mount_fuji:` | [top](#table-of-contents) |
| [top](#travel--places) | :camping: | `:camping:` | :beach_umbrella: | `:beach_umbrella:` | [top](#table-of-contents) |
| [top](#travel--places) | :desert: | `:desert:` | :desert_island: | `:desert_island:` | [top](#table-of-contents) |
| [top](#travel--places) | :national_park: | `:national_park:` | | | [top](#table-of-contents) |

### Place Building

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :stadium: | `:stadium:` | :classical_building: | `:classical_building:` | [top](#table-of-contents) |
| [top](#travel--places) | :building_construction: | `:building_construction:` | :bricks: | `:bricks:` | [top](#table-of-contents) |
| [top](#travel--places) | :rock: | `:rock:` | :wood: | `:wood:` | [top](#table-of-contents) |
| [top](#travel--places) | :hut: | `:hut:` | :houses: | `:houses:` | [top](#table-of-contents) |
| [top](#travel--places) | :derelict_house: | `:derelict_house:` | :house: | `:house:` | [top](#table-of-contents) |
| [top](#travel--places) | :house_with_garden: | `:house_with_garden:` | :office: | `:office:` | [top](#table-of-contents) |
| [top](#travel--places) | :post_office: | `:post_office:` | :european_post_office: | `:european_post_office:` | [top](#table-of-contents) |
| [top](#travel--places) | :hospital: | `:hospital:` | :bank: | `:bank:` | [top](#table-of-contents) |
| [top](#travel--places) | :hotel: | `:hotel:` | :love_hotel: | `:love_hotel:` | [top](#table-of-contents) |
| [top](#travel--places) | :convenience_store: | `:convenience_store:` | :school: | `:school:` | [top](#table-of-contents) |
| [top](#travel--places) | :department_store: | `:department_store:` | :factory: | `:factory:` | [top](#table-of-contents) |
| [top](#travel--places) | :japanese_castle: | `:japanese_castle:` | :european_castle: | `:european_castle:` | [top](#table-of-contents) |
| [top](#travel--places) | :wedding: | `:wedding:` | :tokyo_tower: | `:tokyo_tower:` | [top](#table-of-contents) |
| [top](#travel--places) | :statue_of_liberty: | `:statue_of_liberty:` | | | [top](#table-of-contents) |

### Place Religious

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :church: | `:church:` | :mosque: | `:mosque:` | [top](#table-of-contents) |
| [top](#travel--places) | :hindu_temple: | `:hindu_temple:` | :synagogue: | `:synagogue:` | [top](#table-of-contents) |
| [top](#travel--places) | :shinto_shrine: | `:shinto_shrine:` | :kaaba: | `:kaaba:` | [top](#table-of-contents) |

### Place Other

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :fountain: | `:fountain:` | :tent: | `:tent:` | [top](#table-of-contents) |
| [top](#travel--places) | :foggy: | `:foggy:` | :night_with_stars: | `:night_with_stars:` | [top](#table-of-contents) |
| [top](#travel--places) | :cityscape: | `:cityscape:` | :sunrise_over_mountains: | `:sunrise_over_mountains:` | [top](#table-of-contents) |
| [top](#travel--places) | :sunrise: | `:sunrise:` | :city_sunset: | `:city_sunset:` | [top](#table-of-contents) |
| [top](#travel--places) | :city_sunrise: | `:city_sunrise:` | :bridge_at_night: | `:bridge_at_night:` | [top](#table-of-contents) |
| [top](#travel--places) | :hotsprings: | `:hotsprings:` | :carousel_horse: | `:carousel_horse:` | [top](#table-of-contents) |
| [top](#travel--places) | :playground_slide: | `:playground_slide:` | :ferris_wheel: | `:ferris_wheel:` | [top](#table-of-contents) |
| [top](#travel--places) | :roller_coaster: | `:roller_coaster:` | :barber: | `:barber:` | [top](#table-of-contents) |
| [top](#travel--places) | :circus_tent: | `:circus_tent:` | | | [top](#table-of-contents) |

### Transport Ground

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :steam_locomotive: | `:steam_locomotive:` | :railway_car: | `:railway_car:` | [top](#table-of-contents) |
| [top](#travel--places) | :bullettrain_side: | `:bullettrain_side:` | :bullettrain_front: | `:bullettrain_front:` | [top](#table-of-contents) |
| [top](#travel--places) | :train2: | `:train2:` | :metro: | `:metro:` | [top](#table-of-contents) |
| [top](#travel--places) | :light_rail: | `:light_rail:` | :station: | `:station:` | [top](#table-of-contents) |
| [top](#travel--places) | :tram: | `:tram:` | :monorail: | `:monorail:` | [top](#table-of-contents) |
| [top](#travel--places) | :mountain_railway: | `:mountain_railway:` | :train: | `:train:` | [top](#table-of-contents) |
| [top](#travel--places) | :bus: | `:bus:` | :oncoming_bus: | `:oncoming_bus:` | [top](#table-of-contents) |
| [top](#travel--places) | :trolleybus: | `:trolleybus:` | :minibus: | `:minibus:` | [top](#table-of-contents) |
| [top](#travel--places) | :ambulance: | `:ambulance:` | :fire_engine: | `:fire_engine:` | [top](#table-of-contents) |
| [top](#travel--places) | :police_car: | `:police_car:` | :oncoming_police_car: | `:oncoming_police_car:` | [top](#table-of-contents) |
| [top](#travel--places) | :taxi: | `:taxi:` | :oncoming_taxi: | `:oncoming_taxi:` | [top](#table-of-contents) |
| [top](#travel--places) | :car: | `:car:` <br /> `:red_car:` | :oncoming_automobile: | `:oncoming_automobile:` | [top](#table-of-contents) |
| [top](#travel--places) | :blue_car: | `:blue_car:` | :pickup_truck: | `:pickup_truck:` | [top](#table-of-contents) |
| [top](#travel--places) | :truck: | `:truck:` | :articulated_lorry: | `:articulated_lorry:` | [top](#table-of-contents) |
| [top](#travel--places) | :tractor: | `:tractor:` | :racing_car: | `:racing_car:` | [top](#table-of-contents) |
| [top](#travel--places) | :motorcycle: | `:motorcycle:` | :motor_scooter: | `:motor_scooter:` | [top](#table-of-contents) |
| [top](#travel--places) | :manual_wheelchair: | `:manual_wheelchair:` | :motorized_wheelchair: | `:motorized_wheelchair:` | [top](#table-of-contents) |
| [top](#travel--places) | :auto_rickshaw: | `:auto_rickshaw:` | :bike: | `:bike:` | [top](#table-of-contents) |
| [top](#travel--places) | :kick_scooter: | `:kick_scooter:` | :skateboard: | `:skateboard:` | [top](#table-of-contents) |
| [top](#travel--places) | :roller_skate: | `:roller_skate:` | :busstop: | `:busstop:` | [top](#table-of-contents) |
| [top](#travel--places) | :motorway: | `:motorway:` | :railway_track: | `:railway_track:` | [top](#table-of-contents) |
| [top](#travel--places) | :oil_drum: | `:oil_drum:` | :fuelpump: | `:fuelpump:` | [top](#table-of-contents) |
| [top](#travel--places) | :wheel: | `:wheel:` | :rotating_light: | `:rotating_light:` | [top](#table-of-contents) |
| [top](#travel--places) | :traffic_light: | `:traffic_light:` | :vertical_traffic_light: | `:vertical_traffic_light:` | [top](#table-of-contents) |
| [top](#travel--places) | :stop_sign: | `:stop_sign:` | :construction: | `:construction:` | [top](#table-of-contents) |

### Transport Water

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :anchor: | `:anchor:` | :ring_buoy: | `:ring_buoy:` | [top](#table-of-contents) |
| [top](#travel--places) | :boat: | `:boat:` <br /> `:sailboat:` | :canoe: | `:canoe:` | [top](#table-of-contents) |
| [top](#travel--places) | :speedboat: | `:speedboat:` | :passenger_ship: | `:passenger_ship:` | [top](#table-of-contents) |
| [top](#travel--places) | :ferry: | `:ferry:` | :motor_boat: | `:motor_boat:` | [top](#table-of-contents) |
| [top](#travel--places) | :ship: | `:ship:` | | | [top](#table-of-contents) |

### Transport Air

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :airplane: | `:airplane:` | :small_airplane: | `:small_airplane:` | [top](#table-of-contents) |
| [top](#travel--places) | :flight_departure: | `:flight_departure:` | :flight_arrival: | `:flight_arrival:` | [top](#table-of-contents) |
| [top](#travel--places) | :parachute: | `:parachute:` | :seat: | `:seat:` | [top](#table-of-contents) |
| [top](#travel--places) | :helicopter: | `:helicopter:` | :suspension_railway: | `:suspension_railway:` | [top](#table-of-contents) |
| [top](#travel--places) | :mountain_cableway: | `:mountain_cableway:` | :aerial_tramway: | `:aerial_tramway:` | [top](#table-of-contents) |
| [top](#travel--places) | :artificial_satellite: | `:artificial_satellite:` | :rocket: | `:rocket:` | [top](#table-of-contents) |
| [top](#travel--places) | :flying_saucer: | `:flying_saucer:` | | | [top](#table-of-contents) |

### Hotel

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :bellhop_bell: | `:bellhop_bell:` | :luggage: | `:luggage:` | [top](#table-of-contents) |

### Time

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :hourglass: | `:hourglass:` | :hourglass_flowing_sand: | `:hourglass_flowing_sand:` | [top](#table-of-contents) |
| [top](#travel--places) | :watch: | `:watch:` | :alarm_clock: | `:alarm_clock:` | [top](#table-of-contents) |
| [top](#travel--places) | :stopwatch: | `:stopwatch:` | :timer_clock: | `:timer_clock:` | [top](#table-of-contents) |
| [top](#travel--places) | :mantelpiece_clock: | `:mantelpiece_clock:` | :clock12: | `:clock12:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock1230: | `:clock1230:` | :clock1: | `:clock1:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock130: | `:clock130:` | :clock2: | `:clock2:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock230: | `:clock230:` | :clock3: | `:clock3:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock330: | `:clock330:` | :clock4: | `:clock4:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock430: | `:clock430:` | :clock5: | `:clock5:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock530: | `:clock530:` | :clock6: | `:clock6:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock630: | `:clock630:` | :clock7: | `:clock7:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock730: | `:clock730:` | :clock8: | `:clock8:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock830: | `:clock830:` | :clock9: | `:clock9:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock930: | `:clock930:` | :clock10: | `:clock10:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock1030: | `:clock1030:` | :clock11: | `:clock11:` | [top](#table-of-contents) |
| [top](#travel--places) | :clock1130: | `:clock1130:` | | | [top](#table-of-contents) |

### Sky & Weather

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#travel--places) | :new_moon: | `:new_moon:` | :waxing_crescent_moon: | `:waxing_crescent_moon:` | [top](#table-of-contents) |
| [top](#travel--places) | :first_quarter_moon: | `:first_quarter_moon:` | :moon: | `:moon:` <br /> `:waxing_gibbous_moon:` | [top](#table-of-contents) |
| [top](#travel--places) | :full_moon: | `:full_moon:` | :waning_gibbous_moon: | `:waning_gibbous_moon:` | [top](#table-of-contents) |
| [top](#travel--places) | :last_quarter_moon: | `:last_quarter_moon:` | :waning_crescent_moon: | `:waning_crescent_moon:` | [top](#table-of-contents) |
| [top](#travel--places) | :crescent_moon: | `:crescent_moon:` | :new_moon_with_face: | `:new_moon_with_face:` | [top](#table-of-contents) |
| [top](#travel--places) | :first_quarter_moon_with_face: | `:first_quarter_moon_with_face:` | :last_quarter_moon_with_face: | `:last_quarter_moon_with_face:` | [top](#table-of-contents) |
| [top](#travel--places) | :thermometer: | `:thermometer:` | :sunny: | `:sunny:` | [top](#table-of-contents) |
| [top](#travel--places) | :full_moon_with_face: | `:full_moon_with_face:` | :sun_with_face: | `:sun_with_face:` | [top](#table-of-contents) |
| [top](#travel--places) | :ringed_planet: | `:ringed_planet:` | :star: | `:star:` | [top](#table-of-contents) |
| [top](#travel--places) | :star2: | `:star2:` | :stars: | `:stars:` | [top](#table-of-contents) |
| [top](#travel--places) | :milky_way: | `:milky_way:` | :cloud: | `:cloud:` | [top](#table-of-contents) |
| [top](#travel--places) | :partly_sunny: | `:partly_sunny:` | :cloud_with_lightning_and_rain: | `:cloud_with_lightning_and_rain:` | [top](#table-of-contents) |
| [top](#travel--places) | :sun_behind_small_cloud: | `:sun_behind_small_cloud:` | :sun_behind_large_cloud: | `:sun_behind_large_cloud:` | [top](#table-of-contents) |
| [top](#travel--places) | :sun_behind_rain_cloud: | `:sun_behind_rain_cloud:` | :cloud_with_rain: | `:cloud_with_rain:` | [top](#table-of-contents) |
| [top](#travel--places) | :cloud_with_snow: | `:cloud_with_snow:` | :cloud_with_lightning: | `:cloud_with_lightning:` | [top](#table-of-contents) |
| [top](#travel--places) | :tornado: | `:tornado:` | :fog: | `:fog:` | [top](#table-of-contents) |
| [top](#travel--places) | :wind_face: | `:wind_face:` | :cyclone: | `:cyclone:` | [top](#table-of-contents) |
| [top](#travel--places) | :rainbow: | `:rainbow:` | :closed_umbrella: | `:closed_umbrella:` | [top](#table-of-contents) |
| [top](#travel--places) | :open_umbrella: | `:open_umbrella:` | :umbrella: | `:umbrella:` | [top](#table-of-contents) |
| [top](#travel--places) | :parasol_on_ground: | `:parasol_on_ground:` | :zap: | `:zap:` | [top](#table-of-contents) |
| [top](#travel--places) | :snowflake: | `:snowflake:` | :snowman_with_snow: | `:snowman_with_snow:` | [top](#table-of-contents) |
| [top](#travel--places) | :snowman: | `:snowman:` | :comet: | `:comet:` | [top](#table-of-contents) |
| [top](#travel--places) | :fire: | `:fire:` | :droplet: | `:droplet:` | [top](#table-of-contents) |
| [top](#travel--places) | :ocean: | `:ocean:` | | | [top](#table-of-contents) |

## Activities

- [Event](#event)
- [Award Medal](#award-medal)
- [Sport](#sport)
- [Game](#game)
- [Arts & Crafts](#arts--crafts)

### Event

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#activities) | :jack_o_lantern: | `:jack_o_lantern:` | :christmas_tree: | `:christmas_tree:` | [top](#table-of-contents) |
| [top](#activities) | :fireworks: | `:fireworks:` | :sparkler: | `:sparkler:` | [top](#table-of-contents) |
| [top](#activities) | :firecracker: | `:firecracker:` | :sparkles: | `:sparkles:` | [top](#table-of-contents) |
| [top](#activities) | :balloon: | `:balloon:` | :tada: | `:tada:` | [top](#table-of-contents) |
| [top](#activities) | :confetti_ball: | `:confetti_ball:` | :tanabata_tree: | `:tanabata_tree:` | [top](#table-of-contents) |
| [top](#activities) | :bamboo: | `:bamboo:` | :dolls: | `:dolls:` | [top](#table-of-contents) |
| [top](#activities) | :flags: | `:flags:` | :wind_chime: | `:wind_chime:` | [top](#table-of-contents) |
| [top](#activities) | :rice_scene: | `:rice_scene:` | :red_envelope: | `:red_envelope:` | [top](#table-of-contents) |
| [top](#activities) | :ribbon: | `:ribbon:` | :gift: | `:gift:` | [top](#table-of-contents) |
| [top](#activities) | :reminder_ribbon: | `:reminder_ribbon:` | :tickets: | `:tickets:` | [top](#table-of-contents) |
| [top](#activities) | :ticket: | `:ticket:` | | | [top](#table-of-contents) |

### Award Medal

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#activities) | :medal_military: | `:medal_military:` | :trophy: | `:trophy:` | [top](#table-of-contents) |
| [top](#activities) | :medal_sports: | `:medal_sports:` | :1st_place_medal: | `:1st_place_medal:` | [top](#table-of-contents) |
| [top](#activities) | :2nd_place_medal: | `:2nd_place_medal:` | :3rd_place_medal: | `:3rd_place_medal:` | [top](#table-of-contents) |

### Sport

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#activities) | :soccer: | `:soccer:` | :baseball: | `:baseball:` | [top](#table-of-contents) |
| [top](#activities) | :softball: | `:softball:` | :basketball: | `:basketball:` | [top](#table-of-contents) |
| [top](#activities) | :volleyball: | `:volleyball:` | :football: | `:football:` | [top](#table-of-contents) |
| [top](#activities) | :rugby_football: | `:rugby_football:` | :tennis: | `:tennis:` | [top](#table-of-contents) |
| [top](#activities) | :flying_disc: | `:flying_disc:` | :bowling: | `:bowling:` | [top](#table-of-contents) |
| [top](#activities) | :cricket_game: | `:cricket_game:` | :field_hockey: | `:field_hockey:` | [top](#table-of-contents) |
| [top](#activities) | :ice_hockey: | `:ice_hockey:` | :lacrosse: | `:lacrosse:` | [top](#table-of-contents) |
| [top](#activities) | :ping_pong: | `:ping_pong:` | :badminton: | `:badminton:` | [top](#table-of-contents) |
| [top](#activities) | :boxing_glove: | `:boxing_glove:` | :martial_arts_uniform: | `:martial_arts_uniform:` | [top](#table-of-contents) |
| [top](#activities) | :goal_net: | `:goal_net:` | :golf: | `:golf:` | [top](#table-of-contents) |
| [top](#activities) | :ice_skate: | `:ice_skate:` | :fishing_pole_and_fish: | `:fishing_pole_and_fish:` | [top](#table-of-contents) |
| [top](#activities) | :diving_mask: | `:diving_mask:` | :running_shirt_with_sash: | `:running_shirt_with_sash:` | [top](#table-of-contents) |
| [top](#activities) | :ski: | `:ski:` | :sled: | `:sled:` | [top](#table-of-contents) |
| [top](#activities) | :curling_stone: | `:curling_stone:` | | | [top](#table-of-contents) |

### Game

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#activities) | :dart: | `:dart:` | :yo_yo: | `:yo_yo:` | [top](#table-of-contents) |
| [top](#activities) | :kite: | `:kite:` | :gun: | `:gun:` | [top](#table-of-contents) |
| [top](#activities) | :8ball: | `:8ball:` | :crystal_ball: | `:crystal_ball:` | [top](#table-of-contents) |
| [top](#activities) | :magic_wand: | `:magic_wand:` | :video_game: | `:video_game:` | [top](#table-of-contents) |
| [top](#activities) | :joystick: | `:joystick:` | :slot_machine: | `:slot_machine:` | [top](#table-of-contents) |
| [top](#activities) | :game_die: | `:game_die:` | :jigsaw: | `:jigsaw:` | [top](#table-of-contents) |
| [top](#activities) | :teddy_bear: | `:teddy_bear:` | :pinata: | `:pinata:` | [top](#table-of-contents) |
| [top](#activities) | :mirror_ball: | `:mirror_ball:` | :nesting_dolls: | `:nesting_dolls:` | [top](#table-of-contents) |
| [top](#activities) | :spades: | `:spades:` | :hearts: | `:hearts:` | [top](#table-of-contents) |
| [top](#activities) | :diamonds: | `:diamonds:` | :clubs: | `:clubs:` | [top](#table-of-contents) |
| [top](#activities) | :chess_pawn: | `:chess_pawn:` | :black_joker: | `:black_joker:` | [top](#table-of-contents) |
| [top](#activities) | :mahjong: | `:mahjong:` | :flower_playing_cards: | `:flower_playing_cards:` | [top](#table-of-contents) |

### Arts & Crafts

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#activities) | :performing_arts: | `:performing_arts:` | :framed_picture: | `:framed_picture:` | [top](#table-of-contents) |
| [top](#activities) | :art: | `:art:` | :thread: | `:thread:` | [top](#table-of-contents) |
| [top](#activities) | :sewing_needle: | `:sewing_needle:` | :yarn: | `:yarn:` | [top](#table-of-contents) |
| [top](#activities) | :knot: | `:knot:` | | | [top](#table-of-contents) |

## Objects

- [Clothing](#clothing)
- [Sound](#sound)
- [Music](#music)
- [Musical Instrument](#musical-instrument)
- [Phone](#phone)
- [Computer](#computer)
- [Light & Video](#light--video)
- [Book Paper](#book-paper)
- [Money](#money)
- [Mail](#mail)
- [Writing](#writing)
- [Office](#office)
- [Lock](#lock)
- [Tool](#tool)
- [Science](#science)
- [Medical](#medical)
- [Household](#household)
- [Other Object](#other-object)

### Clothing

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :eyeglasses: | `:eyeglasses:` | :dark_sunglasses: | `:dark_sunglasses:` | [top](#table-of-contents) |
| [top](#objects) | :goggles: | `:goggles:` | :lab_coat: | `:lab_coat:` | [top](#table-of-contents) |
| [top](#objects) | :safety_vest: | `:safety_vest:` | :necktie: | `:necktie:` | [top](#table-of-contents) |
| [top](#objects) | :shirt: | `:shirt:` <br /> `:tshirt:` | :jeans: | `:jeans:` | [top](#table-of-contents) |
| [top](#objects) | :scarf: | `:scarf:` | :gloves: | `:gloves:` | [top](#table-of-contents) |
| [top](#objects) | :coat: | `:coat:` | :socks: | `:socks:` | [top](#table-of-contents) |
| [top](#objects) | :dress: | `:dress:` | :kimono: | `:kimono:` | [top](#table-of-contents) |
| [top](#objects) | :sari: | `:sari:` | :one_piece_swimsuit: | `:one_piece_swimsuit:` | [top](#table-of-contents) |
| [top](#objects) | :swim_brief: | `:swim_brief:` | :shorts: | `:shorts:` | [top](#table-of-contents) |
| [top](#objects) | :bikini: | `:bikini:` | :womans_clothes: | `:womans_clothes:` | [top](#table-of-contents) |
| [top](#objects) | :folding_hand_fan: | `:folding_hand_fan:` | :purse: | `:purse:` | [top](#table-of-contents) |
| [top](#objects) | :handbag: | `:handbag:` | :pouch: | `:pouch:` | [top](#table-of-contents) |
| [top](#objects) | :shopping: | `:shopping:` | :school_satchel: | `:school_satchel:` | [top](#table-of-contents) |
| [top](#objects) | :thong_sandal: | `:thong_sandal:` | :mans_shoe: | `:mans_shoe:` <br /> `:shoe:` | [top](#table-of-contents) |
| [top](#objects) | :athletic_shoe: | `:athletic_shoe:` | :hiking_boot: | `:hiking_boot:` | [top](#table-of-contents) |
| [top](#objects) | :flat_shoe: | `:flat_shoe:` | :high_heel: | `:high_heel:` | [top](#table-of-contents) |
| [top](#objects) | :sandal: | `:sandal:` | :ballet_shoes: | `:ballet_shoes:` | [top](#table-of-contents) |
| [top](#objects) | :boot: | `:boot:` | :hair_pick: | `:hair_pick:` | [top](#table-of-contents) |
| [top](#objects) | :crown: | `:crown:` | :womans_hat: | `:womans_hat:` | [top](#table-of-contents) |
| [top](#objects) | :tophat: | `:tophat:` | :mortar_board: | `:mortar_board:` | [top](#table-of-contents) |
| [top](#objects) | :billed_cap: | `:billed_cap:` | :military_helmet: | `:military_helmet:` | [top](#table-of-contents) |
| [top](#objects) | :rescue_worker_helmet: | `:rescue_worker_helmet:` | :prayer_beads: | `:prayer_beads:` | [top](#table-of-contents) |
| [top](#objects) | :lipstick: | `:lipstick:` | :ring: | `:ring:` | [top](#table-of-contents) |
| [top](#objects) | :gem: | `:gem:` | | | [top](#table-of-contents) |

### Sound

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :mute: | `:mute:` | :speaker: | `:speaker:` | [top](#table-of-contents) |
| [top](#objects) | :sound: | `:sound:` | :loud_sound: | `:loud_sound:` | [top](#table-of-contents) |
| [top](#objects) | :loudspeaker: | `:loudspeaker:` | :mega: | `:mega:` | [top](#table-of-contents) |
| [top](#objects) | :postal_horn: | `:postal_horn:` | :bell: | `:bell:` | [top](#table-of-contents) |
| [top](#objects) | :no_bell: | `:no_bell:` | | | [top](#table-of-contents) |

### Music

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :musical_score: | `:musical_score:` | :musical_note: | `:musical_note:` | [top](#table-of-contents) |
| [top](#objects) | :notes: | `:notes:` | :studio_microphone: | `:studio_microphone:` | [top](#table-of-contents) |
| [top](#objects) | :level_slider: | `:level_slider:` | :control_knobs: | `:control_knobs:` | [top](#table-of-contents) |
| [top](#objects) | :microphone: | `:microphone:` | :headphones: | `:headphones:` | [top](#table-of-contents) |
| [top](#objects) | :radio: | `:radio:` | | | [top](#table-of-contents) |

### Musical Instrument

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :saxophone: | `:saxophone:` | :accordion: | `:accordion:` | [top](#table-of-contents) |
| [top](#objects) | :guitar: | `:guitar:` | :musical_keyboard: | `:musical_keyboard:` | [top](#table-of-contents) |
| [top](#objects) | :trumpet: | `:trumpet:` | :violin: | `:violin:` | [top](#table-of-contents) |
| [top](#objects) | :banjo: | `:banjo:` | :drum: | `:drum:` | [top](#table-of-contents) |
| [top](#objects) | :long_drum: | `:long_drum:` | :maracas: | `:maracas:` | [top](#table-of-contents) |
| [top](#objects) | :flute: | `:flute:` | | | [top](#table-of-contents) |

### Phone

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :iphone: | `:iphone:` | :calling: | `:calling:` | [top](#table-of-contents) |
| [top](#objects) | :phone: | `:phone:` <br /> `:telephone:` | :telephone_receiver: | `:telephone_receiver:` | [top](#table-of-contents) |
| [top](#objects) | :pager: | `:pager:` | :fax: | `:fax:` | [top](#table-of-contents) |

### Computer

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :battery: | `:battery:` | :low_battery: | `:low_battery:` | [top](#table-of-contents) |
| [top](#objects) | :electric_plug: | `:electric_plug:` | :computer: | `:computer:` | [top](#table-of-contents) |
| [top](#objects) | :desktop_computer: | `:desktop_computer:` | :printer: | `:printer:` | [top](#table-of-contents) |
| [top](#objects) | :keyboard: | `:keyboard:` | :computer_mouse: | `:computer_mouse:` | [top](#table-of-contents) |
| [top](#objects) | :trackball: | `:trackball:` | :minidisc: | `:minidisc:` | [top](#table-of-contents) |
| [top](#objects) | :floppy_disk: | `:floppy_disk:` | :cd: | `:cd:` | [top](#table-of-contents) |
| [top](#objects) | :dvd: | `:dvd:` | :abacus: | `:abacus:` | [top](#table-of-contents) |

### Light & Video

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :movie_camera: | `:movie_camera:` | :film_strip: | `:film_strip:` | [top](#table-of-contents) |
| [top](#objects) | :film_projector: | `:film_projector:` | :clapper: | `:clapper:` | [top](#table-of-contents) |
| [top](#objects) | :tv: | `:tv:` | :camera: | `:camera:` | [top](#table-of-contents) |
| [top](#objects) | :camera_flash: | `:camera_flash:` | :video_camera: | `:video_camera:` | [top](#table-of-contents) |
| [top](#objects) | :vhs: | `:vhs:` | :mag: | `:mag:` | [top](#table-of-contents) |
| [top](#objects) | :mag_right: | `:mag_right:` | :candle: | `:candle:` | [top](#table-of-contents) |
| [top](#objects) | :bulb: | `:bulb:` | :flashlight: | `:flashlight:` | [top](#table-of-contents) |
| [top](#objects) | :izakaya_lantern: | `:izakaya_lantern:` <br /> `:lantern:` | :diya_lamp: | `:diya_lamp:` | [top](#table-of-contents) |

### Book Paper

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :notebook_with_decorative_cover: | `:notebook_with_decorative_cover:` | :closed_book: | `:closed_book:` | [top](#table-of-contents) |
| [top](#objects) | :book: | `:book:` <br /> `:open_book:` | :green_book: | `:green_book:` | [top](#table-of-contents) |
| [top](#objects) | :blue_book: | `:blue_book:` | :orange_book: | `:orange_book:` | [top](#table-of-contents) |
| [top](#objects) | :books: | `:books:` | :notebook: | `:notebook:` | [top](#table-of-contents) |
| [top](#objects) | :ledger: | `:ledger:` | :page_with_curl: | `:page_with_curl:` | [top](#table-of-contents) |
| [top](#objects) | :scroll: | `:scroll:` | :page_facing_up: | `:page_facing_up:` | [top](#table-of-contents) |
| [top](#objects) | :newspaper: | `:newspaper:` | :newspaper_roll: | `:newspaper_roll:` | [top](#table-of-contents) |
| [top](#objects) | :bookmark_tabs: | `:bookmark_tabs:` | :bookmark: | `:bookmark:` | [top](#table-of-contents) |
| [top](#objects) | :label: | `:label:` | | | [top](#table-of-contents) |

### Money

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :moneybag: | `:moneybag:` | :coin: | `:coin:` | [top](#table-of-contents) |
| [top](#objects) | :yen: | `:yen:` | :dollar: | `:dollar:` | [top](#table-of-contents) |
| [top](#objects) | :euro: | `:euro:` | :pound: | `:pound:` | [top](#table-of-contents) |
| [top](#objects) | :money_with_wings: | `:money_with_wings:` | :credit_card: | `:credit_card:` | [top](#table-of-contents) |
| [top](#objects) | :receipt: | `:receipt:` | :chart: | `:chart:` | [top](#table-of-contents) |

### Mail

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :envelope: | `:envelope:` | :e-mail: | `:e-mail:` <br /> `:email:` | [top](#table-of-contents) |
| [top](#objects) | :incoming_envelope: | `:incoming_envelope:` | :envelope_with_arrow: | `:envelope_with_arrow:` | [top](#table-of-contents) |
| [top](#objects) | :outbox_tray: | `:outbox_tray:` | :inbox_tray: | `:inbox_tray:` | [top](#table-of-contents) |
| [top](#objects) | :package: | `:package:` | :mailbox: | `:mailbox:` | [top](#table-of-contents) |
| [top](#objects) | :mailbox_closed: | `:mailbox_closed:` | :mailbox_with_mail: | `:mailbox_with_mail:` | [top](#table-of-contents) |
| [top](#objects) | :mailbox_with_no_mail: | `:mailbox_with_no_mail:` | :postbox: | `:postbox:` | [top](#table-of-contents) |
| [top](#objects) | :ballot_box: | `:ballot_box:` | | | [top](#table-of-contents) |

### Writing

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :pencil2: | `:pencil2:` | :black_nib: | `:black_nib:` | [top](#table-of-contents) |
| [top](#objects) | :fountain_pen: | `:fountain_pen:` | :pen: | `:pen:` | [top](#table-of-contents) |
| [top](#objects) | :paintbrush: | `:paintbrush:` | :crayon: | `:crayon:` | [top](#table-of-contents) |
| [top](#objects) | :memo: | `:memo:` <br /> `:pencil:` | | | [top](#table-of-contents) |

### Office

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :briefcase: | `:briefcase:` | :file_folder: | `:file_folder:` | [top](#table-of-contents) |
| [top](#objects) | :open_file_folder: | `:open_file_folder:` | :card_index_dividers: | `:card_index_dividers:` | [top](#table-of-contents) |
| [top](#objects) | :date: | `:date:` | :calendar: | `:calendar:` | [top](#table-of-contents) |
| [top](#objects) | :spiral_notepad: | `:spiral_notepad:` | :spiral_calendar: | `:spiral_calendar:` | [top](#table-of-contents) |
| [top](#objects) | :card_index: | `:card_index:` | :chart_with_upwards_trend: | `:chart_with_upwards_trend:` | [top](#table-of-contents) |
| [top](#objects) | :chart_with_downwards_trend: | `:chart_with_downwards_trend:` | :bar_chart: | `:bar_chart:` | [top](#table-of-contents) |
| [top](#objects) | :clipboard: | `:clipboard:` | :pushpin: | `:pushpin:` | [top](#table-of-contents) |
| [top](#objects) | :round_pushpin: | `:round_pushpin:` | :paperclip: | `:paperclip:` | [top](#table-of-contents) |
| [top](#objects) | :paperclips: | `:paperclips:` | :straight_ruler: | `:straight_ruler:` | [top](#table-of-contents) |
| [top](#objects) | :triangular_ruler: | `:triangular_ruler:` | :scissors: | `:scissors:` | [top](#table-of-contents) |
| [top](#objects) | :card_file_box: | `:card_file_box:` | :file_cabinet: | `:file_cabinet:` | [top](#table-of-contents) |
| [top](#objects) | :wastebasket: | `:wastebasket:` | | | [top](#table-of-contents) |

### Lock

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :lock: | `:lock:` | :unlock: | `:unlock:` | [top](#table-of-contents) |
| [top](#objects) | :lock_with_ink_pen: | `:lock_with_ink_pen:` | :closed_lock_with_key: | `:closed_lock_with_key:` | [top](#table-of-contents) |
| [top](#objects) | :key: | `:key:` | :old_key: | `:old_key:` | [top](#table-of-contents) |

### Tool

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :hammer: | `:hammer:` | :axe: | `:axe:` | [top](#table-of-contents) |
| [top](#objects) | :pick: | `:pick:` | :hammer_and_pick: | `:hammer_and_pick:` | [top](#table-of-contents) |
| [top](#objects) | :hammer_and_wrench: | `:hammer_and_wrench:` | :dagger: | `:dagger:` | [top](#table-of-contents) |
| [top](#objects) | :crossed_swords: | `:crossed_swords:` | :bomb: | `:bomb:` | [top](#table-of-contents) |
| [top](#objects) | :boomerang: | `:boomerang:` | :bow_and_arrow: | `:bow_and_arrow:` | [top](#table-of-contents) |
| [top](#objects) | :shield: | `:shield:` | :carpentry_saw: | `:carpentry_saw:` | [top](#table-of-contents) |
| [top](#objects) | :wrench: | `:wrench:` | :screwdriver: | `:screwdriver:` | [top](#table-of-contents) |
| [top](#objects) | :nut_and_bolt: | `:nut_and_bolt:` | :gear: | `:gear:` | [top](#table-of-contents) |
| [top](#objects) | :clamp: | `:clamp:` | :balance_scale: | `:balance_scale:` | [top](#table-of-contents) |
| [top](#objects) | :probing_cane: | `:probing_cane:` | :link: | `:link:` | [top](#table-of-contents) |
| [top](#objects) | :chains: | `:chains:` | :hook: | `:hook:` | [top](#table-of-contents) |
| [top](#objects) | :toolbox: | `:toolbox:` | :magnet: | `:magnet:` | [top](#table-of-contents) |
| [top](#objects) | :ladder: | `:ladder:` | | | [top](#table-of-contents) |

### Science

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :alembic: | `:alembic:` | :test_tube: | `:test_tube:` | [top](#table-of-contents) |
| [top](#objects) | :petri_dish: | `:petri_dish:` | :dna: | `:dna:` | [top](#table-of-contents) |
| [top](#objects) | :microscope: | `:microscope:` | :telescope: | `:telescope:` | [top](#table-of-contents) |
| [top](#objects) | :satellite: | `:satellite:` | | | [top](#table-of-contents) |

### Medical

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :syringe: | `:syringe:` | :drop_of_blood: | `:drop_of_blood:` | [top](#table-of-contents) |
| [top](#objects) | :pill: | `:pill:` | :adhesive_bandage: | `:adhesive_bandage:` | [top](#table-of-contents) |
| [top](#objects) | :crutch: | `:crutch:` | :stethoscope: | `:stethoscope:` | [top](#table-of-contents) |
| [top](#objects) | :x_ray: | `:x_ray:` | | | [top](#table-of-contents) |

### Household

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :door: | `:door:` | :elevator: | `:elevator:` | [top](#table-of-contents) |
| [top](#objects) | :mirror: | `:mirror:` | :window: | `:window:` | [top](#table-of-contents) |
| [top](#objects) | :bed: | `:bed:` | :couch_and_lamp: | `:couch_and_lamp:` | [top](#table-of-contents) |
| [top](#objects) | :chair: | `:chair:` | :toilet: | `:toilet:` | [top](#table-of-contents) |
| [top](#objects) | :plunger: | `:plunger:` | :shower: | `:shower:` | [top](#table-of-contents) |
| [top](#objects) | :bathtub: | `:bathtub:` | :mouse_trap: | `:mouse_trap:` | [top](#table-of-contents) |
| [top](#objects) | :razor: | `:razor:` | :lotion_bottle: | `:lotion_bottle:` | [top](#table-of-contents) |
| [top](#objects) | :safety_pin: | `:safety_pin:` | :broom: | `:broom:` | [top](#table-of-contents) |
| [top](#objects) | :basket: | `:basket:` | :roll_of_paper: | `:roll_of_paper:` | [top](#table-of-contents) |
| [top](#objects) | :bucket: | `:bucket:` | :soap: | `:soap:` | [top](#table-of-contents) |
| [top](#objects) | :bubbles: | `:bubbles:` | :toothbrush: | `:toothbrush:` | [top](#table-of-contents) |
| [top](#objects) | :sponge: | `:sponge:` | :fire_extinguisher: | `:fire_extinguisher:` | [top](#table-of-contents) |
| [top](#objects) | :shopping_cart: | `:shopping_cart:` | | | [top](#table-of-contents) |

### Other Object

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#objects) | :smoking: | `:smoking:` | :coffin: | `:coffin:` | [top](#table-of-contents) |
| [top](#objects) | :headstone: | `:headstone:` | :funeral_urn: | `:funeral_urn:` | [top](#table-of-contents) |
| [top](#objects) | :nazar_amulet: | `:nazar_amulet:` | :hamsa: | `:hamsa:` | [top](#table-of-contents) |
| [top](#objects) | :moyai: | `:moyai:` | :placard: | `:placard:` | [top](#table-of-contents) |
| [top](#objects) | :identification_card: | `:identification_card:` | | | [top](#table-of-contents) |

## Symbols

- [Transport Sign](#transport-sign)
- [Warning](#warning)
- [Arrow](#arrow)
- [Religion](#religion)
- [Zodiac](#zodiac)
- [Av Symbol](#av-symbol)
- [Gender](#gender)
- [Math](#math)
- [Punctuation](#punctuation)
- [Currency](#currency)
- [Other Symbol](#other-symbol)
- [Keycap](#keycap)
- [Alphanum](#alphanum)
- [Geometric](#geometric)

### Transport Sign

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :atm: | `:atm:` | :put_litter_in_its_place: | `:put_litter_in_its_place:` | [top](#table-of-contents) |
| [top](#symbols) | :potable_water: | `:potable_water:` | :wheelchair: | `:wheelchair:` | [top](#table-of-contents) |
| [top](#symbols) | :mens: | `:mens:` | :womens: | `:womens:` | [top](#table-of-contents) |
| [top](#symbols) | :restroom: | `:restroom:` | :baby_symbol: | `:baby_symbol:` | [top](#table-of-contents) |
| [top](#symbols) | :wc: | `:wc:` | :passport_control: | `:passport_control:` | [top](#table-of-contents) |
| [top](#symbols) | :customs: | `:customs:` | :baggage_claim: | `:baggage_claim:` | [top](#table-of-contents) |
| [top](#symbols) | :left_luggage: | `:left_luggage:` | | | [top](#table-of-contents) |

### Warning

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :warning: | `:warning:` | :children_crossing: | `:children_crossing:` | [top](#table-of-contents) |
| [top](#symbols) | :no_entry: | `:no_entry:` | :no_entry_sign: | `:no_entry_sign:` | [top](#table-of-contents) |
| [top](#symbols) | :no_bicycles: | `:no_bicycles:` | :no_smoking: | `:no_smoking:` | [top](#table-of-contents) |
| [top](#symbols) | :do_not_litter: | `:do_not_litter:` | :non-potable_water: | `:non-potable_water:` | [top](#table-of-contents) |
| [top](#symbols) | :no_pedestrians: | `:no_pedestrians:` | :no_mobile_phones: | `:no_mobile_phones:` | [top](#table-of-contents) |
| [top](#symbols) | :underage: | `:underage:` | :radioactive: | `:radioactive:` | [top](#table-of-contents) |
| [top](#symbols) | :biohazard: | `:biohazard:` | | | [top](#table-of-contents) |

### Arrow

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :arrow_up: | `:arrow_up:` | :arrow_upper_right: | `:arrow_upper_right:` | [top](#table-of-contents) |
| [top](#symbols) | :arrow_right: | `:arrow_right:` | :arrow_lower_right: | `:arrow_lower_right:` | [top](#table-of-contents) |
| [top](#symbols) | :arrow_down: | `:arrow_down:` | :arrow_lower_left: | `:arrow_lower_left:` | [top](#table-of-contents) |
| [top](#symbols) | :arrow_left: | `:arrow_left:` | :arrow_upper_left: | `:arrow_upper_left:` | [top](#table-of-contents) |
| [top](#symbols) | :arrow_up_down: | `:arrow_up_down:` | :left_right_arrow: | `:left_right_arrow:` | [top](#table-of-contents) |
| [top](#symbols) | :leftwards_arrow_with_hook: | `:leftwards_arrow_with_hook:` | :arrow_right_hook: | `:arrow_right_hook:` | [top](#table-of-contents) |
| [top](#symbols) | :arrow_heading_up: | `:arrow_heading_up:` | :arrow_heading_down: | `:arrow_heading_down:` | [top](#table-of-contents) |
| [top](#symbols) | :arrows_clockwise: | `:arrows_clockwise:` | :arrows_counterclockwise: | `:arrows_counterclockwise:` | [top](#table-of-contents) |
| [top](#symbols) | :back: | `:back:` | :end: | `:end:` | [top](#table-of-contents) |
| [top](#symbols) | :on: | `:on:` | :soon: | `:soon:` | [top](#table-of-contents) |
| [top](#symbols) | :top: | `:top:` | | | [top](#table-of-contents) |

### Religion

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :place_of_worship: | `:place_of_worship:` | :atom_symbol: | `:atom_symbol:` | [top](#table-of-contents) |
| [top](#symbols) | :om: | `:om:` | :star_of_david: | `:star_of_david:` | [top](#table-of-contents) |
| [top](#symbols) | :wheel_of_dharma: | `:wheel_of_dharma:` | :yin_yang: | `:yin_yang:` | [top](#table-of-contents) |
| [top](#symbols) | :latin_cross: | `:latin_cross:` | :orthodox_cross: | `:orthodox_cross:` | [top](#table-of-contents) |
| [top](#symbols) | :star_and_crescent: | `:star_and_crescent:` | :peace_symbol: | `:peace_symbol:` | [top](#table-of-contents) |
| [top](#symbols) | :menorah: | `:menorah:` | :six_pointed_star: | `:six_pointed_star:` | [top](#table-of-contents) |
| [top](#symbols) | :khanda: | `:khanda:` | | | [top](#table-of-contents) |

### Zodiac

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :aries: | `:aries:` | :taurus: | `:taurus:` | [top](#table-of-contents) |
| [top](#symbols) | :gemini: | `:gemini:` | :cancer: | `:cancer:` | [top](#table-of-contents) |
| [top](#symbols) | :leo: | `:leo:` | :virgo: | `:virgo:` | [top](#table-of-contents) |
| [top](#symbols) | :libra: | `:libra:` | :scorpius: | `:scorpius:` | [top](#table-of-contents) |
| [top](#symbols) | :sagittarius: | `:sagittarius:` | :capricorn: | `:capricorn:` | [top](#table-of-contents) |
| [top](#symbols) | :aquarius: | `:aquarius:` | :pisces: | `:pisces:` | [top](#table-of-contents) |
| [top](#symbols) | :ophiuchus: | `:ophiuchus:` | | | [top](#table-of-contents) |

### Av Symbol

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :twisted_rightwards_arrows: | `:twisted_rightwards_arrows:` | :repeat: | `:repeat:` | [top](#table-of-contents) |
| [top](#symbols) | :repeat_one: | `:repeat_one:` | :arrow_forward: | `:arrow_forward:` | [top](#table-of-contents) |
| [top](#symbols) | :fast_forward: | `:fast_forward:` | :next_track_button: | `:next_track_button:` | [top](#table-of-contents) |
| [top](#symbols) | :play_or_pause_button: | `:play_or_pause_button:` | :arrow_backward: | `:arrow_backward:` | [top](#table-of-contents) |
| [top](#symbols) | :rewind: | `:rewind:` | :previous_track_button: | `:previous_track_button:` | [top](#table-of-contents) |
| [top](#symbols) | :arrow_up_small: | `:arrow_up_small:` | :arrow_double_up: | `:arrow_double_up:` | [top](#table-of-contents) |
| [top](#symbols) | :arrow_down_small: | `:arrow_down_small:` | :arrow_double_down: | `:arrow_double_down:` | [top](#table-of-contents) |
| [top](#symbols) | :pause_button: | `:pause_button:` | :stop_button: | `:stop_button:` | [top](#table-of-contents) |
| [top](#symbols) | :record_button: | `:record_button:` | :eject_button: | `:eject_button:` | [top](#table-of-contents) |
| [top](#symbols) | :cinema: | `:cinema:` | :low_brightness: | `:low_brightness:` | [top](#table-of-contents) |
| [top](#symbols) | :high_brightness: | `:high_brightness:` | :signal_strength: | `:signal_strength:` | [top](#table-of-contents) |
| [top](#symbols) | :wireless: | `:wireless:` | :vibration_mode: | `:vibration_mode:` | [top](#table-of-contents) |
| [top](#symbols) | :mobile_phone_off: | `:mobile_phone_off:` | | | [top](#table-of-contents) |

### Gender

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :female_sign: | `:female_sign:` | :male_sign: | `:male_sign:` | [top](#table-of-contents) |
| [top](#symbols) | :transgender_symbol: | `:transgender_symbol:` | | | [top](#table-of-contents) |

### Math

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :heavy_multiplication_x: | `:heavy_multiplication_x:` | :heavy_plus_sign: | `:heavy_plus_sign:` | [top](#table-of-contents) |
| [top](#symbols) | :heavy_minus_sign: | `:heavy_minus_sign:` | :heavy_division_sign: | `:heavy_division_sign:` | [top](#table-of-contents) |
| [top](#symbols) | :heavy_equals_sign: | `:heavy_equals_sign:` | :infinity: | `:infinity:` | [top](#table-of-contents) |

### Punctuation

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :bangbang: | `:bangbang:` | :interrobang: | `:interrobang:` | [top](#table-of-contents) |
| [top](#symbols) | :question: | `:question:` | :grey_question: | `:grey_question:` | [top](#table-of-contents) |
| [top](#symbols) | :grey_exclamation: | `:grey_exclamation:` | :exclamation: | `:exclamation:` <br /> `:heavy_exclamation_mark:` | [top](#table-of-contents) |
| [top](#symbols) | :wavy_dash: | `:wavy_dash:` | | | [top](#table-of-contents) |

### Currency

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :currency_exchange: | `:currency_exchange:` | :heavy_dollar_sign: | `:heavy_dollar_sign:` | [top](#table-of-contents) |

### Other Symbol

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :medical_symbol: | `:medical_symbol:` | :recycle: | `:recycle:` | [top](#table-of-contents) |
| [top](#symbols) | :fleur_de_lis: | `:fleur_de_lis:` | :trident: | `:trident:` | [top](#table-of-contents) |
| [top](#symbols) | :name_badge: | `:name_badge:` | :beginner: | `:beginner:` | [top](#table-of-contents) |
| [top](#symbols) | :o: | `:o:` | :white_check_mark: | `:white_check_mark:` | [top](#table-of-contents) |
| [top](#symbols) | :ballot_box_with_check: | `:ballot_box_with_check:` | :heavy_check_mark: | `:heavy_check_mark:` | [top](#table-of-contents) |
| [top](#symbols) | :x: | `:x:` | :negative_squared_cross_mark: | `:negative_squared_cross_mark:` | [top](#table-of-contents) |
| [top](#symbols) | :curly_loop: | `:curly_loop:` | :loop: | `:loop:` | [top](#table-of-contents) |
| [top](#symbols) | :part_alternation_mark: | `:part_alternation_mark:` | :eight_spoked_asterisk: | `:eight_spoked_asterisk:` | [top](#table-of-contents) |
| [top](#symbols) | :eight_pointed_black_star: | `:eight_pointed_black_star:` | :sparkle: | `:sparkle:` | [top](#table-of-contents) |
| [top](#symbols) | :copyright: | `:copyright:` | :registered: | `:registered:` | [top](#table-of-contents) |
| [top](#symbols) | :tm: | `:tm:` | | | [top](#table-of-contents) |

### Keycap

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :hash: | `:hash:` | :asterisk: | `:asterisk:` | [top](#table-of-contents) |
| [top](#symbols) | :zero: | `:zero:` | :one: | `:one:` | [top](#table-of-contents) |
| [top](#symbols) | :two: | `:two:` | :three: | `:three:` | [top](#table-of-contents) |
| [top](#symbols) | :four: | `:four:` | :five: | `:five:` | [top](#table-of-contents) |
| [top](#symbols) | :six: | `:six:` | :seven: | `:seven:` | [top](#table-of-contents) |
| [top](#symbols) | :eight: | `:eight:` | :nine: | `:nine:` | [top](#table-of-contents) |
| [top](#symbols) | :keycap_ten: | `:keycap_ten:` | | | [top](#table-of-contents) |

### Alphanum

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :capital_abcd: | `:capital_abcd:` | :abcd: | `:abcd:` | [top](#table-of-contents) |
| [top](#symbols) | :1234: | `:1234:` | :symbols: | `:symbols:` | [top](#table-of-contents) |
| [top](#symbols) | :abc: | `:abc:` | :a: | `:a:` | [top](#table-of-contents) |
| [top](#symbols) | :ab: | `:ab:` | :b: | `:b:` | [top](#table-of-contents) |
| [top](#symbols) | :cl: | `:cl:` | :cool: | `:cool:` | [top](#table-of-contents) |
| [top](#symbols) | :free: | `:free:` | :information_source: | `:information_source:` | [top](#table-of-contents) |
| [top](#symbols) | :id: | `:id:` | :m: | `:m:` | [top](#table-of-contents) |
| [top](#symbols) | :new: | `:new:` | :ng: | `:ng:` | [top](#table-of-contents) |
| [top](#symbols) | :o2: | `:o2:` | :ok: | `:ok:` | [top](#table-of-contents) |
| [top](#symbols) | :parking: | `:parking:` | :sos: | `:sos:` | [top](#table-of-contents) |
| [top](#symbols) | :up: | `:up:` | :vs: | `:vs:` | [top](#table-of-contents) |
| [top](#symbols) | :koko: | `:koko:` | :sa: | `:sa:` | [top](#table-of-contents) |
| [top](#symbols) | :u6708: | `:u6708:` | :u6709: | `:u6709:` | [top](#table-of-contents) |
| [top](#symbols) | :u6307: | `:u6307:` | :ideograph_advantage: | `:ideograph_advantage:` | [top](#table-of-contents) |
| [top](#symbols) | :u5272: | `:u5272:` | :u7121: | `:u7121:` | [top](#table-of-contents) |
| [top](#symbols) | :u7981: | `:u7981:` | :accept: | `:accept:` | [top](#table-of-contents) |
| [top](#symbols) | :u7533: | `:u7533:` | :u5408: | `:u5408:` | [top](#table-of-contents) |
| [top](#symbols) | :u7a7a: | `:u7a7a:` | :congratulations: | `:congratulations:` | [top](#table-of-contents) |
| [top](#symbols) | :secret: | `:secret:` | :u55b6: | `:u55b6:` | [top](#table-of-contents) |
| [top](#symbols) | :u6e80: | `:u6e80:` | | | [top](#table-of-contents) |

### Geometric

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#symbols) | :red_circle: | `:red_circle:` | :orange_circle: | `:orange_circle:` | [top](#table-of-contents) |
| [top](#symbols) | :yellow_circle: | `:yellow_circle:` | :green_circle: | `:green_circle:` | [top](#table-of-contents) |
| [top](#symbols) | :large_blue_circle: | `:large_blue_circle:` | :purple_circle: | `:purple_circle:` | [top](#table-of-contents) |
| [top](#symbols) | :brown_circle: | `:brown_circle:` | :black_circle: | `:black_circle:` | [top](#table-of-contents) |
| [top](#symbols) | :white_circle: | `:white_circle:` | :red_square: | `:red_square:` | [top](#table-of-contents) |
| [top](#symbols) | :orange_square: | `:orange_square:` | :yellow_square: | `:yellow_square:` | [top](#table-of-contents) |
| [top](#symbols) | :green_square: | `:green_square:` | :blue_square: | `:blue_square:` | [top](#table-of-contents) |
| [top](#symbols) | :purple_square: | `:purple_square:` | :brown_square: | `:brown_square:` | [top](#table-of-contents) |
| [top](#symbols) | :black_large_square: | `:black_large_square:` | :white_large_square: | `:white_large_square:` | [top](#table-of-contents) |
| [top](#symbols) | :black_medium_square: | `:black_medium_square:` | :white_medium_square: | `:white_medium_square:` | [top](#table-of-contents) |
| [top](#symbols) | :black_medium_small_square: | `:black_medium_small_square:` | :white_medium_small_square: | `:white_medium_small_square:` | [top](#table-of-contents) |
| [top](#symbols) | :black_small_square: | `:black_small_square:` | :white_small_square: | `:white_small_square:` | [top](#table-of-contents) |
| [top](#symbols) | :large_orange_diamond: | `:large_orange_diamond:` | :large_blue_diamond: | `:large_blue_diamond:` | [top](#table-of-contents) |
| [top](#symbols) | :small_orange_diamond: | `:small_orange_diamond:` | :small_blue_diamond: | `:small_blue_diamond:` | [top](#table-of-contents) |
| [top](#symbols) | :small_red_triangle: | `:small_red_triangle:` | :small_red_triangle_down: | `:small_red_triangle_down:` | [top](#table-of-contents) |
| [top](#symbols) | :diamond_shape_with_a_dot_inside: | `:diamond_shape_with_a_dot_inside:` | :radio_button: | `:radio_button:` | [top](#table-of-contents) |
| [top](#symbols) | :white_square_button: | `:white_square_button:` | :black_square_button: | `:black_square_button:` | [top](#table-of-contents) |

## Flags

- [Flag](#flag)
- [Country Flag](#country-flag)
- [Subdivision Flag](#subdivision-flag)

### Flag

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#flags) | :checkered_flag: | `:checkered_flag:` | :triangular_flag_on_post: | `:triangular_flag_on_post:` | [top](#table-of-contents) |
| [top](#flags) | :crossed_flags: | `:crossed_flags:` | :black_flag: | `:black_flag:` | [top](#table-of-contents) |
| [top](#flags) | :white_flag: | `:white_flag:` | :rainbow_flag: | `:rainbow_flag:` | [top](#table-of-contents) |
| [top](#flags) | :transgender_flag: | `:transgender_flag:` | :pirate_flag: | `:pirate_flag:` | [top](#table-of-contents) |

### Country Flag

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#flags) | :ascension_island: | `:ascension_island:` | :andorra: | `:andorra:` | [top](#table-of-contents) |
| [top](#flags) | :united_arab_emirates: | `:united_arab_emirates:` | :afghanistan: | `:afghanistan:` | [top](#table-of-contents) |
| [top](#flags) | :antigua_barbuda: | `:antigua_barbuda:` | :anguilla: | `:anguilla:` | [top](#table-of-contents) |
| [top](#flags) | :albania: | `:albania:` | :armenia: | `:armenia:` | [top](#table-of-contents) |
| [top](#flags) | :angola: | `:angola:` | :antarctica: | `:antarctica:` | [top](#table-of-contents) |
| [top](#flags) | :argentina: | `:argentina:` | :american_samoa: | `:american_samoa:` | [top](#table-of-contents) |
| [top](#flags) | :austria: | `:austria:` | :australia: | `:australia:` | [top](#table-of-contents) |
| [top](#flags) | :aruba: | `:aruba:` | :aland_islands: | `:aland_islands:` | [top](#table-of-contents) |
| [top](#flags) | :azerbaijan: | `:azerbaijan:` | :bosnia_herzegovina: | `:bosnia_herzegovina:` | [top](#table-of-contents) |
| [top](#flags) | :barbados: | `:barbados:` | :bangladesh: | `:bangladesh:` | [top](#table-of-contents) |
| [top](#flags) | :belgium: | `:belgium:` | :burkina_faso: | `:burkina_faso:` | [top](#table-of-contents) |
| [top](#flags) | :bulgaria: | `:bulgaria:` | :bahrain: | `:bahrain:` | [top](#table-of-contents) |
| [top](#flags) | :burundi: | `:burundi:` | :benin: | `:benin:` | [top](#table-of-contents) |
| [top](#flags) | :st_barthelemy: | `:st_barthelemy:` | :bermuda: | `:bermuda:` | [top](#table-of-contents) |
| [top](#flags) | :brunei: | `:brunei:` | :bolivia: | `:bolivia:` | [top](#table-of-contents) |
| [top](#flags) | :caribbean_netherlands: | `:caribbean_netherlands:` | :brazil: | `:brazil:` | [top](#table-of-contents) |
| [top](#flags) | :bahamas: | `:bahamas:` | :bhutan: | `:bhutan:` | [top](#table-of-contents) |
| [top](#flags) | :bouvet_island: | `:bouvet_island:` | :botswana: | `:botswana:` | [top](#table-of-contents) |
| [top](#flags) | :belarus: | `:belarus:` | :belize: | `:belize:` | [top](#table-of-contents) |
| [top](#flags) | :canada: | `:canada:` | :cocos_islands: | `:cocos_islands:` | [top](#table-of-contents) |
| [top](#flags) | :congo_kinshasa: | `:congo_kinshasa:` | :central_african_republic: | `:central_african_republic:` | [top](#table-of-contents) |
| [top](#flags) | :congo_brazzaville: | `:congo_brazzaville:` | :switzerland: | `:switzerland:` | [top](#table-of-contents) |
| [top](#flags) | :cote_divoire: | `:cote_divoire:` | :cook_islands: | `:cook_islands:` | [top](#table-of-contents) |
| [top](#flags) | :chile: | `:chile:` | :cameroon: | `:cameroon:` | [top](#table-of-contents) |
| [top](#flags) | :cn: | `:cn:` | :colombia: | `:colombia:` | [top](#table-of-contents) |
| [top](#flags) | :clipperton_island: | `:clipperton_island:` | :costa_rica: | `:costa_rica:` | [top](#table-of-contents) |
| [top](#flags) | :cuba: | `:cuba:` | :cape_verde: | `:cape_verde:` | [top](#table-of-contents) |
| [top](#flags) | :curacao: | `:curacao:` | :christmas_island: | `:christmas_island:` | [top](#table-of-contents) |
| [top](#flags) | :cyprus: | `:cyprus:` | :czech_republic: | `:czech_republic:` | [top](#table-of-contents) |
| [top](#flags) | :de: | `:de:` | :diego_garcia: | `:diego_garcia:` | [top](#table-of-contents) |
| [top](#flags) | :djibouti: | `:djibouti:` | :denmark: | `:denmark:` | [top](#table-of-contents) |
| [top](#flags) | :dominica: | `:dominica:` | :dominican_republic: | `:dominican_republic:` | [top](#table-of-contents) |
| [top](#flags) | :algeria: | `:algeria:` | :ceuta_melilla: | `:ceuta_melilla:` | [top](#table-of-contents) |
| [top](#flags) | :ecuador: | `:ecuador:` | :estonia: | `:estonia:` | [top](#table-of-contents) |
| [top](#flags) | :egypt: | `:egypt:` | :western_sahara: | `:western_sahara:` | [top](#table-of-contents) |
| [top](#flags) | :eritrea: | `:eritrea:` | :es: | `:es:` | [top](#table-of-contents) |
| [top](#flags) | :ethiopia: | `:ethiopia:` | :eu: | `:eu:` <br /> `:european_union:` | [top](#table-of-contents) |
| [top](#flags) | :finland: | `:finland:` | :fiji: | `:fiji:` | [top](#table-of-contents) |
| [top](#flags) | :falkland_islands: | `:falkland_islands:` | :micronesia: | `:micronesia:` | [top](#table-of-contents) |
| [top](#flags) | :faroe_islands: | `:faroe_islands:` | :fr: | `:fr:` | [top](#table-of-contents) |
| [top](#flags) | :gabon: | `:gabon:` | :gb: | `:gb:` <br /> `:uk:` | [top](#table-of-contents) |
| [top](#flags) | :grenada: | `:grenada:` | :georgia: | `:georgia:` | [top](#table-of-contents) |
| [top](#flags) | :french_guiana: | `:french_guiana:` | :guernsey: | `:guernsey:` | [top](#table-of-contents) |
| [top](#flags) | :ghana: | `:ghana:` | :gibraltar: | `:gibraltar:` | [top](#table-of-contents) |
| [top](#flags) | :greenland: | `:greenland:` | :gambia: | `:gambia:` | [top](#table-of-contents) |
| [top](#flags) | :guinea: | `:guinea:` | :guadeloupe: | `:guadeloupe:` | [top](#table-of-contents) |
| [top](#flags) | :equatorial_guinea: | `:equatorial_guinea:` | :greece: | `:greece:` | [top](#table-of-contents) |
| [top](#flags) | :south_georgia_south_sandwich_islands: | `:south_georgia_south_sandwich_islands:` | :guatemala: | `:guatemala:` | [top](#table-of-contents) |
| [top](#flags) | :guam: | `:guam:` | :guinea_bissau: | `:guinea_bissau:` | [top](#table-of-contents) |
| [top](#flags) | :guyana: | `:guyana:` | :hong_kong: | `:hong_kong:` | [top](#table-of-contents) |
| [top](#flags) | :heard_mcdonald_islands: | `:heard_mcdonald_islands:` | :honduras: | `:honduras:` | [top](#table-of-contents) |
| [top](#flags) | :croatia: | `:croatia:` | :haiti: | `:haiti:` | [top](#table-of-contents) |
| [top](#flags) | :hungary: | `:hungary:` | :canary_islands: | `:canary_islands:` | [top](#table-of-contents) |
| [top](#flags) | :indonesia: | `:indonesia:` | :ireland: | `:ireland:` | [top](#table-of-contents) |
| [top](#flags) | :israel: | `:israel:` | :isle_of_man: | `:isle_of_man:` | [top](#table-of-contents) |
| [top](#flags) | :india: | `:india:` | :british_indian_ocean_territory: | `:british_indian_ocean_territory:` | [top](#table-of-contents) |
| [top](#flags) | :iraq: | `:iraq:` | :iran: | `:iran:` | [top](#table-of-contents) |
| [top](#flags) | :iceland: | `:iceland:` | :it: | `:it:` | [top](#table-of-contents) |
| [top](#flags) | :jersey: | `:jersey:` | :jamaica: | `:jamaica:` | [top](#table-of-contents) |
| [top](#flags) | :jordan: | `:jordan:` | :jp: | `:jp:` | [top](#table-of-contents) |
| [top](#flags) | :kenya: | `:kenya:` | :kyrgyzstan: | `:kyrgyzstan:` | [top](#table-of-contents) |
| [top](#flags) | :cambodia: | `:cambodia:` | :kiribati: | `:kiribati:` | [top](#table-of-contents) |
| [top](#flags) | :comoros: | `:comoros:` | :st_kitts_nevis: | `:st_kitts_nevis:` | [top](#table-of-contents) |
| [top](#flags) | :north_korea: | `:north_korea:` | :kr: | `:kr:` | [top](#table-of-contents) |
| [top](#flags) | :kuwait: | `:kuwait:` | :cayman_islands: | `:cayman_islands:` | [top](#table-of-contents) |
| [top](#flags) | :kazakhstan: | `:kazakhstan:` | :laos: | `:laos:` | [top](#table-of-contents) |
| [top](#flags) | :lebanon: | `:lebanon:` | :st_lucia: | `:st_lucia:` | [top](#table-of-contents) |
| [top](#flags) | :liechtenstein: | `:liechtenstein:` | :sri_lanka: | `:sri_lanka:` | [top](#table-of-contents) |
| [top](#flags) | :liberia: | `:liberia:` | :lesotho: | `:lesotho:` | [top](#table-of-contents) |
| [top](#flags) | :lithuania: | `:lithuania:` | :luxembourg: | `:luxembourg:` | [top](#table-of-contents) |
| [top](#flags) | :latvia: | `:latvia:` | :libya: | `:libya:` | [top](#table-of-contents) |
| [top](#flags) | :morocco: | `:morocco:` | :monaco: | `:monaco:` | [top](#table-of-contents) |
| [top](#flags) | :moldova: | `:moldova:` | :montenegro: | `:montenegro:` | [top](#table-of-contents) |
| [top](#flags) | :st_martin: | `:st_martin:` | :madagascar: | `:madagascar:` | [top](#table-of-contents) |
| [top](#flags) | :marshall_islands: | `:marshall_islands:` | :macedonia: | `:macedonia:` | [top](#table-of-contents) |
| [top](#flags) | :mali: | `:mali:` | :myanmar: | `:myanmar:` | [top](#table-of-contents) |
| [top](#flags) | :mongolia: | `:mongolia:` | :macau: | `:macau:` | [top](#table-of-contents) |
| [top](#flags) | :northern_mariana_islands: | `:northern_mariana_islands:` | :martinique: | `:martinique:` | [top](#table-of-contents) |
| [top](#flags) | :mauritania: | `:mauritania:` | :montserrat: | `:montserrat:` | [top](#table-of-contents) |
| [top](#flags) | :malta: | `:malta:` | :mauritius: | `:mauritius:` | [top](#table-of-contents) |
| [top](#flags) | :maldives: | `:maldives:` | :malawi: | `:malawi:` | [top](#table-of-contents) |
| [top](#flags) | :mexico: | `:mexico:` | :malaysia: | `:malaysia:` | [top](#table-of-contents) |
| [top](#flags) | :mozambique: | `:mozambique:` | :namibia: | `:namibia:` | [top](#table-of-contents) |
| [top](#flags) | :new_caledonia: | `:new_caledonia:` | :niger: | `:niger:` | [top](#table-of-contents) |
| [top](#flags) | :norfolk_island: | `:norfolk_island:` | :nigeria: | `:nigeria:` | [top](#table-of-contents) |
| [top](#flags) | :nicaragua: | `:nicaragua:` | :netherlands: | `:netherlands:` | [top](#table-of-contents) |
| [top](#flags) | :norway: | `:norway:` | :nepal: | `:nepal:` | [top](#table-of-contents) |
| [top](#flags) | :nauru: | `:nauru:` | :niue: | `:niue:` | [top](#table-of-contents) |
| [top](#flags) | :new_zealand: | `:new_zealand:` | :oman: | `:oman:` | [top](#table-of-contents) |
| [top](#flags) | :panama: | `:panama:` | :peru: | `:peru:` | [top](#table-of-contents) |
| [top](#flags) | :french_polynesia: | `:french_polynesia:` | :papua_new_guinea: | `:papua_new_guinea:` | [top](#table-of-contents) |
| [top](#flags) | :philippines: | `:philippines:` | :pakistan: | `:pakistan:` | [top](#table-of-contents) |
| [top](#flags) | :poland: | `:poland:` | :st_pierre_miquelon: | `:st_pierre_miquelon:` | [top](#table-of-contents) |
| [top](#flags) | :pitcairn_islands: | `:pitcairn_islands:` | :puerto_rico: | `:puerto_rico:` | [top](#table-of-contents) |
| [top](#flags) | :palestinian_territories: | `:palestinian_territories:` | :portugal: | `:portugal:` | [top](#table-of-contents) |
| [top](#flags) | :palau: | `:palau:` | :paraguay: | `:paraguay:` | [top](#table-of-contents) |
| [top](#flags) | :qatar: | `:qatar:` | :reunion: | `:reunion:` | [top](#table-of-contents) |
| [top](#flags) | :romania: | `:romania:` | :serbia: | `:serbia:` | [top](#table-of-contents) |
| [top](#flags) | :ru: | `:ru:` | :rwanda: | `:rwanda:` | [top](#table-of-contents) |
| [top](#flags) | :saudi_arabia: | `:saudi_arabia:` | :solomon_islands: | `:solomon_islands:` | [top](#table-of-contents) |
| [top](#flags) | :seychelles: | `:seychelles:` | :sudan: | `:sudan:` | [top](#table-of-contents) |
| [top](#flags) | :sweden: | `:sweden:` | :singapore: | `:singapore:` | [top](#table-of-contents) |
| [top](#flags) | :st_helena: | `:st_helena:` | :slovenia: | `:slovenia:` | [top](#table-of-contents) |
| [top](#flags) | :svalbard_jan_mayen: | `:svalbard_jan_mayen:` | :slovakia: | `:slovakia:` | [top](#table-of-contents) |
| [top](#flags) | :sierra_leone: | `:sierra_leone:` | :san_marino: | `:san_marino:` | [top](#table-of-contents) |
| [top](#flags) | :senegal: | `:senegal:` | :somalia: | `:somalia:` | [top](#table-of-contents) |
| [top](#flags) | :suriname: | `:suriname:` | :south_sudan: | `:south_sudan:` | [top](#table-of-contents) |
| [top](#flags) | :sao_tome_principe: | `:sao_tome_principe:` | :el_salvador: | `:el_salvador:` | [top](#table-of-contents) |
| [top](#flags) | :sint_maarten: | `:sint_maarten:` | :syria: | `:syria:` | [top](#table-of-contents) |
| [top](#flags) | :swaziland: | `:swaziland:` | :tristan_da_cunha: | `:tristan_da_cunha:` | [top](#table-of-contents) |
| [top](#flags) | :turks_caicos_islands: | `:turks_caicos_islands:` | :chad: | `:chad:` | [top](#table-of-contents) |
| [top](#flags) | :french_southern_territories: | `:french_southern_territories:` | :togo: | `:togo:` | [top](#table-of-contents) |
| [top](#flags) | :thailand: | `:thailand:` | :tajikistan: | `:tajikistan:` | [top](#table-of-contents) |
| [top](#flags) | :tokelau: | `:tokelau:` | :timor_leste: | `:timor_leste:` | [top](#table-of-contents) |
| [top](#flags) | :turkmenistan: | `:turkmenistan:` | :tunisia: | `:tunisia:` | [top](#table-of-contents) |
| [top](#flags) | :tonga: | `:tonga:` | :tr: | `:tr:` | [top](#table-of-contents) |
| [top](#flags) | :trinidad_tobago: | `:trinidad_tobago:` | :tuvalu: | `:tuvalu:` | [top](#table-of-contents) |
| [top](#flags) | :taiwan: | `:taiwan:` | :tanzania: | `:tanzania:` | [top](#table-of-contents) |
| [top](#flags) | :ukraine: | `:ukraine:` | :uganda: | `:uganda:` | [top](#table-of-contents) |
| [top](#flags) | :us_outlying_islands: | `:us_outlying_islands:` | :united_nations: | `:united_nations:` | [top](#table-of-contents) |
| [top](#flags) | :us: | `:us:` | :uruguay: | `:uruguay:` | [top](#table-of-contents) |
| [top](#flags) | :uzbekistan: | `:uzbekistan:` | :vatican_city: | `:vatican_city:` | [top](#table-of-contents) |
| [top](#flags) | :st_vincent_grenadines: | `:st_vincent_grenadines:` | :venezuela: | `:venezuela:` | [top](#table-of-contents) |
| [top](#flags) | :british_virgin_islands: | `:british_virgin_islands:` | :us_virgin_islands: | `:us_virgin_islands:` | [top](#table-of-contents) |
| [top](#flags) | :vietnam: | `:vietnam:` | :vanuatu: | `:vanuatu:` | [top](#table-of-contents) |
| [top](#flags) | :wallis_futuna: | `:wallis_futuna:` | :samoa: | `:samoa:` | [top](#table-of-contents) |
| [top](#flags) | :kosovo: | `:kosovo:` | :yemen: | `:yemen:` | [top](#table-of-contents) |
| [top](#flags) | :mayotte: | `:mayotte:` | :south_africa: | `:south_africa:` | [top](#table-of-contents) |
| [top](#flags) | :zambia: | `:zambia:` | :zimbabwe: | `:zimbabwe:` | [top](#table-of-contents) |

### Subdivision Flag

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#flags) | :england: | `:england:` | :scotland: | `:scotland:` | [top](#table-of-contents) |
| [top](#flags) | :wales: | `:wales:` | | | [top](#table-of-contents) |

## GitHub Custom Emoji

| | ico | shortcode | ico | shortcode | |
| - | :-: | - | :-: | - | - |
| [top](#github-custom-emoji) | :accessibility: | `:accessibility:` | :atom: | `:atom:` | [top](#table-of-contents) |
| [top](#github-custom-emoji) | :basecamp: | `:basecamp:` | :basecampy: | `:basecampy:` | [top](#table-of-contents) |
| [top](#github-custom-emoji) | :bowtie: | `:bowtie:` | :dependabot: | `:dependabot:` | [top](#table-of-contents) |
| [top](#github-custom-emoji) | :electron: | `:electron:` | :feelsgood: | `:feelsgood:` | [top](#table-of-contents) |
| [top](#github-custom-emoji) | :finnadie: | `:finnadie:` | :fishsticks: | `:fishsticks:` | [top](#table-of-contents) |
| [top](#github-custom-emoji) | :goberserk: | `:goberserk:` | :godmode: | `:godmode:` | [top](#table-of-contents) |
| [top](#github-custom-emoji) | :hurtrealbad: | `:hurtrealbad:` | :neckbeard: | `:neckbeard:` | [top](#table-of-contents) |
| [top](#github-custom-emoji) | :octocat: | `:octocat:` | :rage1: | `:rage1:` | [top](#table-of-contents) |
| [top](#github-custom-emoji) | :rage2: | `:rage2:` | :rage3: | `:rage3:` | [top](#table-of-contents) |
| [top](#github-custom-emoji) | :rage4: | `:rage4:` | :shipit: | `:shipit:` | [top](#table-of-contents) |
| [top](#github-custom-emoji) | :suspect: | `:suspect:` | :trollface: | `:trollface:` | [top](#table-of-contents) |
