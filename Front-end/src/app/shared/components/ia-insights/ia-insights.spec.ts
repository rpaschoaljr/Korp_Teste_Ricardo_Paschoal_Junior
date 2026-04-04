import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IaInsights } from './ia-insights';

describe('IaInsights', () => {
  let component: IaInsights;
  let fixture: ComponentFixture<IaInsights>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [IaInsights]
    })
    .compileComponents();

    fixture = TestBed.createComponent(IaInsights);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
